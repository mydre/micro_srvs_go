package main

import (
	"go.uber.org/zap"
)
func main(){

	logger,_ := zap.NewProduction()//生产环境
	//开发环境NewDevelopment
	defer logger.Sync()
	url := "https://imooc.com"
	logger.Info("failed to fetch URL",
		zap.String("url",url),zap.Int("nums",3),)
	//{"level":"info","ts":1625552886.5805845,"caller":"zap_test/main.go:12","msg":"failed to fetch URL","url":"https://imooc.com","nums":3}
	//sugar := logger.Sugar()
	///*
	//	Sugarreddit Logger 和 Logger的对比，Logger比Sugar Logger快，但是Sugar更加方便
	//*/
	//sugar.Infow("failed to fetch URL", "url",url,"attempt",3,)
	//sugar.Infof("Failed to fetch URL: %s,",url)
}