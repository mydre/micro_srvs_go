package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"
)

var (
	ServerConfig *config.ServerConfig= &config.ServerConfig{}
	Trans        ut.Translator

	UserSrvClient proto.UserClient

	NacosConfig *config.NacosConfig = &config.NacosConfig{}//viper做的工作是初始化nacos config
)
