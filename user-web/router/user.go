package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func InitUserRouter(router *gin.RouterGroup){
	//UserGroup := router.Group("user").Use(middlewares.JWTAuth())
	UserGroup := router.Group("user").Use()
	zap.S().Info("配置用户相关的URL")
	{
		UserGroup.GET("list", middlewares.JWTAuth(),middlewares.IsAdminAuth(),api.GetUserList)
		UserGroup.POST("pwd_login", api.PassWordLogin)//这是post method
		UserGroup.POST("register",api.Register)
	}
}