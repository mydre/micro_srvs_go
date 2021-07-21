package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func InitSrvConn(){
	zap.S().Info("初始化Client连接")
	consuleInfo := global.ServerConfig.ConsulInfo //获取注册中心的ip 和 port
	//"consul://10.114.21.16:8500/user-srv?wait=14s&tag=imooc"
	connStr := fmt.Sprintf("consul://%s:%d/%s?wait=14s&tag=imooc",consuleInfo.Host,consuleInfo.Port,global.ServerConfig.UserSrvInfo.Name)
	conn, err := grpc.Dial(//之前拨号连接的是service层，现在拨号连接的是注册中心
		//"consul://10.114.21.16:8500/user-srv?wait=14s&tag=imooc",
		connStr,
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
	}
	userSrvClient := proto.NewUserClient(conn)
	global.UserSrvClient = userSrvClient//这里就一直建立了连接，而不会关闭
}

func InitSrvConn2(){
	zap.S().Info("初始化Client")
	cfg := api.DefaultConfig()
	consuleInfo := global.ServerConfig.ConsulInfo //获取注册中心的ip 和 port
	cfg.Address = fmt.Sprintf("%s:%d",consuleInfo.Host,consuleInfo.Port)
	userSrvHost := ""
	userSrvPort := 0

	client,err := api.NewClient(cfg)
	if err!=nil{
		panic(err)
	}
	//获取python的grpc提供的service
	data,err :=client.Agent().ServicesWithFilter(fmt.Sprintf(`Service=="%s"`,global.ServerConfig.UserSrvInfo.Name))

	if err != nil{
		panic(err)
	}
	for _,value := range data{
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == ""{
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}
	/*
			浏览器在什么情况下会发起options预检请求？
			在非简单请求(设置beforeSend请求)且跨域的情况下，浏览器会发起options预检请求。
		简单请求：请求方法是HEAD，POST，GET三者之一，且http的header中不超出以下字段：accept,accept-language,
		content-language,last-event-id,content-type
	*/
	//ip := global.ServerConfig.UserSrvInfo.Host
	//port := global.ServerConfig.UserSrvInfo.Port
	//拨号连接用户grpc服务
	userConn,err := grpc.Dial(fmt.Sprintf("%s:%d",userSrvHost,userSrvPort),grpc.WithInsecure())
	if err!=nil{
		zap.S().Errorw("[GetUserList] 连接 [用户服务失败]",
			"msg",err.Error(),)
	}
	//1.后续的用户服务下线了？全局变量怎样改？2.改端口了 3.改ip了， 都会导致出问题，
	//2.已经事先创立好了连接，这样后续就不用进行再次tcp的三次握手了
	//3.一个连接多个groutine（协程）共用，性能 - 连接池(连接池有开源的项目)
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
}
