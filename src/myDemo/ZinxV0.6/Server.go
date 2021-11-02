package main

import (
	"awesomeProject/src/zinx/ziface"
	"awesomeProject/src/zinx/znet"
	"fmt"
)
// inherit base Router
type PingRouter struct {
	znet.BaseRouter
}
// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ... ")
	// read from client-side, and ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil{
		fmt.Println(err)
	}
}

// 根据不同的消息处理不同的业务

// hello Zinx test customized router
type HelloZinxRouter struct{
	znet.BaseRouter
}
// Test Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ... ")
	// read from client-side, and ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hellow Welcome to Zinx!"))
	if err != nil{
		fmt.Println(err)
	}
}





func main() {
	//1. create server handler, use Zinx Api
	s := znet.NewServer("[zinxv0.6]")


	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})


	// customize AddRouter
	//s.AddRouter(&PingRouter{})

	//2. initalize server
	s.Serve()

}
