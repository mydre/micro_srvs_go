package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/utils"
)

func GetEnvInfo(env string)bool{
	viper.AutomaticEnv()
	return viper.GetBool(env)

}
func InitConfig(){
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml",configFilePrefix)
	if debug{
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml",configFilePrefix)
	}
	v := viper.New()
	//文件路径
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil{
		panic(err)
	}
	//serverConfig := *global.ServerConfig,serverConfig是一个变量，存放其的地址是另一个地址，和global.ServerConfig是不同的
	if err := v.Unmarshal(global.NacosConfig); err != nil{//viper能够将yaml配置文件的信息映射到global.ServerConfig中.通过指针进行访问
		panic(err)
	}
	zap.S().Infof("配置信息：%v",global.NacosConfig)

	//viper的功能 - 动态监控变化(viper监听本地文件的时候不会阻塞程序的运行)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event){
		zap.S().Infof("配置文件发生变化：%s",e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.NacosConfig)
		zap.S().Infof("配置信息：&v",global.NacosConfig)
	})
	//从nacos中读取配置信息
	utils.NacosConfigSupervise(global.NacosConfig)
}