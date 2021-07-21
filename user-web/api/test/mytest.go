package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
	"time"
)

func loadBanlance(){//负载均衡相关的测试代码
	conn, err := grpc.Dial(//之前拨号连接的是service层，现在拨号连接的是注册中心
		"consul://10.114.21.16:8500/user-srv?wait=14s&tag=imooc",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		//log.Fatal(err)
		log.Fatal(err)
	}
	defer conn.Close()
	for i:=0;i<10;i++{
		userSrvClient := proto.NewUserClient(conn)
		zap.S().Info(context.Background())
		rsp,err := userSrvClient.GetUserList(context.Background(),&proto.PageInfo{
			Pn:1,
			PSize: 2,
		})
		if err != nil{
			panic(err)
		}
		for index,data := range rsp.Data{
			fmt.Println(index,data)
		}
	}
}

func nacosConfigSupervise(){
	//nacosConfig是指向NacosConfig的指针

	sc := []constant.ServerConfig{
		{
			IpAddr: "10.114.21.16",
			Port: 8848,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         "65b3b123-1a2b-4053-8802-7665aa3356e2", //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
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
		DataId: "user-web.json",
		Group: "dev",
	})
	if err !=  nil{
		panic(err)
	}
	/*
		将nacos监听的json配置文件对应的config,转换为go中的struct对象
		想要将json字符串转换为go中的struct对象,需要在struct对象中设置tag
	*/
	//其中content是字符串:string类型
	json.Unmarshal([]byte(content),global.ServerConfig)//global.ServerConfig是一个指针

	fmt.Printf("content:\n",content)
	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-web.json",
		Group: "dev",
		OnChange: func(namespace,group,dataId,data string){
			fmt.Println("配置文件变化")
			fmt.Println("namespace:" + namespace + ", group:" + group + ", dataId:" + dataId + ", data:\n" + data)
		},
	})
	time.Sleep(3000 * time.Second)
}

func main() {
}