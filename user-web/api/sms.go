package api

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"net/http"
	"strings"
	"time"
)

func GenerateSmsCode(width int)string{
	//生成width长度的短信验证码
	numeric := [10]byte{0,1,2,3,4,5,6,7,8,9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0;i<width;i++{
		fmt.Fprintf(&sb,"%d",numeric[rand.Intn(r)])
	}
	//return sb.String()
	return "123456"//由于没有短信服务这里把验证码写固定，假设用户知道验证码是123456
}

func SendSms(ctx *gin.Context){//其中会调用GenerateSmsCode
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm);err != nil{
		HandleValidatorError(ctx,err)
		return
	}
	client,err := dysmsapi.NewClientWithAccessKey("cn-beijing",global.ServerConfig.AliSmsInfo.ApiKey,global.ServerConfig.AliSmsInfo.ApiSecret)
	if err!=nil{
		panic(err)
	}

	mobile := sendSmsForm.Mobile
	smsCode := GenerateSmsCode(6)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https"//https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-beijing"
	request.QueryParams["PhoneNumbers"] = mobile
	request.QueryParams["SignName"] = "xxx"
	request.QueryParams["TemplateCode"] = "xxx"
	request.QueryParams["TemplateParam"] = "{\"code\":" + smsCode + "}" //
	response, err := client.ProcessCommonRequest(request)//把随机的验证码发送到用户的手机？

	fmt.Println(client.DoAction(request,response))
	//fmt.Print(response)
	if err!= nil{
		fmt.Println(err.Error())
	}
	fmt.Printf("response is %#v\n",response)//验证短信验证码是否成功
	//业务逻辑,短信发送完成之后，后面注册的时候，会把短信验证码带回来,动态生成制定长度的验证码
	//将验证码和手机号码保存起来,手机号码作为key，验证码作为value

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",global.ServerConfig.RedisInfo.Host,global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(),mobile,smsCode,time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)//设置15秒过期
	//可以给前端返回
	ctx.JSON(http.StatusOK,gin.H{
		"msg":"短信验证码发送成功",
	})
}