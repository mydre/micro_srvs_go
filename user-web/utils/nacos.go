package utils

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/global"
	"time"
)

func NacosConfigSupervise(nacosConfig *config.NacosConfig){
	//nacosConfig是指向NacosConfig的指针

	sc := []constant.ServerConfig{
		{
			IpAddr: nacosConfig.Host,
			Port: uint64(nacosConfig.Port),
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         nacosConfig.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs":  sc,
		"clientConfig": cc,
	})
	if err != nil{
		panic(err)
	}
	content ,err:= configClient.GetConfig(vo.ConfigParam{
		DataId: nacosConfig.DataId,
		Group: nacosConfig.Group,
	})
	if err !=  nil{
		panic(err)
	}
	/*
		将nacos监听的json配置文件对应的config,转换为go中的struct对象
		想要将json字符串转换为go中的struct对象,需要在struct对象中设置tag
	*/
	//其中content是字符串:string类型
	err = json.Unmarshal([]byte(content),global.ServerConfig)//global.ServerConfig是一个指针
	if err != nil{
		zap.S().Fatalf("读取nacos配置失败: %s",err.Error())
	}

	fmt.Printf("content:\n",content)
	zap.S().Info(global.ServerConfig)
	go func(){//放在协程中处理,可以防止阻塞
		err = configClient.ListenConfig(vo.ConfigParam{
			DataId: "user-web.json",
			Group: "dev",
			OnChange: func(namespace,group,dataId,data string){
				fmt.Println("配置文件变化")
				fmt.Println("namespace:" + namespace + ", group:" + group + ", dataId:" + dataId + ", data:\n" + data)
			},
		})
		time.Sleep(3000 * time.Second)
	}()
}