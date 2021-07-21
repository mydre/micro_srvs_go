package utils

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)
/*
	将web服务注册到注册中心
*/
func Register(address string,port int,name string,tags []string, id string)error{
	cfg := api.DefaultConfig()
	cfg.Address = "10.114.21.16:8500"

	client,err := api.NewClient(cfg)
	if err != nil{
		panic(err)
	}

	//new()函数返回的是指向这个参数的指针
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP: "http://10.114.21.16:8021/health",//这里必须是10.114.21.16
		Timeout: "5s",
		Interval: "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	registration.Check = check//registration的check是一个指针
	err = client.Agent().ServiceRegister(registration)
	if err != nil{
		panic(err)
	}
	return nil
	//生成对应的检查对象
}

func AllServices(){
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"

	client,err := api.NewClient(cfg)
	if err != nil{
		panic(err)
	}

	data,err:=client.Agent().Services()//返回值是map，key是string，value是指针
	if err!=nil{
		panic(err)
	}
	for key,_:= range data{//遍历map
		fmt.Println(key)//user-web
	}
}
//拿到感兴趣的服务
func FilterService() {
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := client.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		panic(err)
	}
	for key, _ := range data { //遍历map
		fmt.Println(key) //user-web
	}
}