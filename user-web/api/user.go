package api

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)
var trans ut.Translator

func HandleGrpcErrorToHttp(err error,c *gin.Context){
	//将grpc的code装换成http的状态码
	if err != nil{
		if e,ok := status.FromError(err); ok{
			switch e.Code(){
			case codes.NotFound:
				c.JSON(http.StatusNotFound,gin.H{
					"msg":e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError,gin.H{
					"msg":"内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest,gin.H{
					"msg":"参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError,gin.H{
					"msg":"用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError,gin.H{
					//"msg":"其他错误",
					"msg":e.Code(),
				})
			}
			return
		}
	}
}

func HandleValidatorError(c *gin.Context,err error){
	errs, ok := err.(validator.ValidationErrors)
	if !ok {//说明验证失败
		c.JSON(http.StatusOK, gin.H{//以json的形式返回response
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{//验证
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func GetUserList(ctx *gin.Context){//函数的参数是一个指向Context类型的指针
	claims,_ := ctx.Get("claims")//由于在前面已经加上了权限验证的逻辑，并附带这存入了claims
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户：%d",currentUser.ID)
	userSrvClient := global.UserSrvClient

	//ctx中保存的是请求的参数（可能是post请求，也可能是get请求的参数？）
	pn := ctx.DefaultQuery("pn","0")
	pnInt,_ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize","10")
	pSizeInt,_ := strconv.Atoi(pSize)

	rsp,err := userSrvClient.GetUserList(context.Background(),&proto.PageInfo{
		Pn:uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err!=nil{
		zap.S().Errorw("[GetUserList] 查询 [用户列表] 失败")
		HandleGrpcErrorToHttp(err,ctx)
		return
	}
	//返回数据，可能有多条数据，这个时候，返回JSON
	//result := make(map[string]interface{})
	result := make([]interface{},0)
	for _,value := range rsp.Data{//在go中可以直接检索protobuf的字段！

		//data := make(map[string]interface{})
		user := response.UserResponse{
			Id:value.Id,
			NickName: value.NickName,
			//BirthDay: time.Time(time.Unix(int64(value.BirthDay),0)),
			BirthDay: time.Time(time.Unix(int64(value.BirthDay),0)).Format("2006-01-02"),
			Gender: value.Gender,
			Mobile: value.Mobile,
		}
		result = append(result,user)
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile
		//result = append(result,data)
	}
	//ctx.JSON(http.StatusOK,result)//result是一个json
	ctx.JSON(http.StatusOK,result)//result是一个json
}

func PassWordLogin(c *gin.Context){
	//表单验证
	passwordLoginForm := forms.PassWordLoginForm{}//表单是用来接收并传递值的
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(c,err)
		return
	}

	//首先把图片验证码环节去除
	//if !store.Verify(passwordLoginForm.CaptchaId,passwordLoginForm.Captcha,false){
	//	c.JSON(http.StatusBadRequest,gin.H{
	//		"captcha":"验证码错误",
	//	})
	//	return
	//}


	userSrvClient := global.UserSrvClient
	//登录的逻辑
	zap.S().Info(passwordLoginForm.Mobile)
	if rsp,err := userSrvClient.GetUserByMobile(context.Background(),&proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	});err !=nil{
		if e,ok := status.FromError(err); ok {//查找具体是哪个错误
			switch e.Code(){
			case codes.NotFound:
				c.JSON(http.StatusBadRequest,map[string]string{
					"mobile":"用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError,map[string]string{
					"mobile":"登录【查询用户失败】",
					"code":fmt.Sprintf("%d",e.Code()),
				})
			}
			return//有错误要记得返回
		}
	}else{
		//只是查询到了用户而已，并没有检查密码(web层没有能力自己做密码的检查)
		if passRsp,passErr := userSrvClient.CheckPassWord(context.Background(),&proto.PasswordCheckInfo{
			Password: passwordLoginForm.Password,
			EncryptedPassword:rsp.PassWord,
		});passErr!= nil{
			c.JSON(http.StatusInternalServerError,map[string]string{
				"password":"登录失败",
			})
		}else{
			if passRsp.Success{
				//生成token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID: uint(rsp.Id),
					NickName: rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims:jwt.StandardClaims{
						NotBefore: time.Now().Unix(),//签名的生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30,//30天过期
						Issuer: "imooc",
					},
				}
				token,err := j.CreateToken(claims)
				if err != nil{
					c.JSON(http.StatusInternalServerError,gin.H{
						"msg":"生成token失败",
					})
					return
				}
				c.JSON(http.StatusOK,gin.H{
					"id":rsp.Id,
					"nick_name":rsp.NickName,
					"token":token,
					"expired_at":(time.Now().Unix() + 60*60*24*30)*1000,
				})

				//c.JSON(http.StatusOK,map[string]string{
				//	"msg":"登录成功",
				//})
			}else{
				c.JSON(http.StatusBadRequest,map[string]string{
					"msg":"登录失败,密码错误",
				})
			}

		}
	}
}

func Register(c *gin.Context){
	//表单验证
	registerForm:= forms.RegisterForm{}//表单是用来接收并传递值的
	if err := c.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(c,err)
		return
	}
	//验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",global.ServerConfig.RedisInfo.Host,global.ServerConfig.RedisInfo.Port),
	})
	if value,err := rdb.Get(context.Background(),registerForm.Mobile).Result();err ==redis.Nil{
		//redis.Nil是不存在key的情况
		//fmt.Println("key 不存在")
		c.JSON(http.StatusBadRequest,gin.H{
			"code":"验证码已经过期或还未获得验证码",
		})
		return
	}else{
		if value != registerForm.Code{
			c.JSON(http.StatusBadRequest,gin.H{
				"code":"验证码输入错误",
			})
			return
		}
	}
	//用户注册
	//拨号连接用户grpc服务
	userSrvClient := global.UserSrvClient
	user,err := userSrvClient.CreateUser(context.Background(),&proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.Password,
		Mobile: registerForm.Mobile,
	})
	if err != nil{
		zap.S().Errorf("[Register]【新建用户失败】%s",err.Error())
		HandleGrpcErrorToHttp(err,c)
	}

	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID: uint(user.Id),
		NickName: user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims:jwt.StandardClaims{
			NotBefore: time.Now().Unix(),//签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30,//30天过期
			Issuer: "imooc",
		},
	}
	token,err := j.CreateToken(claims)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"msg":"生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"id":user.Id,
		"nick_name":user.NickName,
		"token":token,
		"expired_at":(time.Now().Unix() + 60*60*24*30)*1000,
	})
}