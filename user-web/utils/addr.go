package utils

import (
	"net"
)

func GetFreePort()(int,error){
	addr,err := net.ResolveTCPAddr("tcp","localhost:0")

	if err != nil{
		return 0,err
	}

	//l是指向TcpListener的指针,通过指针来调用该指针指向的对象上所绑定的函数
	l,err := net.ListenTCP("tcp",addr)
	if err != nil{
		return 0,err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port,nil
}