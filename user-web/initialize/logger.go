package initialize

import "go.uber.org/zap"
/*
	S()可以获取一个全局的sugar,可以让我们设置一个全局的logger
	S函数和L函数很有用,提供了一个全局的安全访问logger的途径
*/
func InitLogger(){
	//logger,_ := zap.NewProduction()//这个是json格式,生产环境的日志级别是INFO，所以不查看DEBUG类型的日志
	logger,_ := zap.NewDevelopment()//这个是日志格式,开发环境的日志级别是DEBUG
	zap.ReplaceGlobals(logger)
}
