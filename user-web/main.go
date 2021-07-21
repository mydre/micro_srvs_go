package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"
	myvalidator "mxshop-api/user-web/validator"
)

func main(){
	viper.AutomaticEnv()
	//如果是本地开发环境端口号固定，线上环境自动获取端口号
	debug := viper.GetBool("MXSHOP_DEBUG")//通过在编辑器直接运行，而不是在终端运行，则获取不到对应的环境变量，则运行与release模式
	if !debug{
		gin.SetMode(gin.ReleaseMode)
	}

	//1.初始化logger
	initialize.InitLogger()
	//2.初始化配置文件
	initialize.InitConfig()
	//3.初始化routers
	Router := initialize.Routers()
	//4.初始化翻译（验证表单的时候变为中文的提示）
	if err := initialize.InitTrans("zh");err !=nil{
		panic(err)
	}
	//5.初始化连接(现在之后srv（用户连接）)
	initialize.InitSrvConn()
	//6.获取自身的端口
	if !debug{
		port ,err:= utils.GetFreePort()
		if err == nil{
			global.ServerConfig.Port = port
		}
	}
	//注册自定义验证器（手机号码验证器）
	if v,ok := binding.Validator.Engine().(*validator.Validate);ok{
		_ = v.RegisterValidation("mobile",myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile",global.Trans,func(ut ut.Translator)error{
			return ut.Add("mobile","{0} 非法的手机号码!",true)
		},func(ut ut.Translator,fe validator.FieldError)string{
			t,_ := ut.T("mobile",fe.Field())
			return t
		})
	}
	//port, _ := strconv.Atoi(global.ServerConfig.Port)//字符串转换为int
	port := global.ServerConfig.Port//字符串转换为int

	//注册服务
	zap.S().Info("注册服务")
	utils.Register("10.114.21.16",port,"user-web",[]string{"mxshop"},"user-web")
	zap.S().Infof("启动服务器,端口：%d",port)
	if err := Router.Run(fmt.Sprintf(":%d",port));err != nil{
		zap.S().Panic("启动失败:",err.Error())
	}
}