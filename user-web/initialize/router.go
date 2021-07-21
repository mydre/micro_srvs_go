package initialize

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/router"
	"net/http"
)

func Routers() *gin.Engine{
	zap.S().Info("初始化路由")
	Router := gin.Default()
	Router.GET("/health",func(c *gin.Context){
		c.JSON(http.StatusOK,gin.H{
			"code":http.StatusOK,
			"success":true,
		})
	})
	//配置跨域
	Router.Use(middlewares.Cors())
	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return Router
}