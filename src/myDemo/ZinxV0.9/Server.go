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

	// 发送TLV格式的
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

//  execute Hook after create conn,
func DOConnectionBegin(conn ziface.IConnection){
	fmt.Println("==> DoConnectionBegin is Called ... ")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil{
		fmt.Println(err)
	}
}

// execute after disconnect
func DoConnectionLost(conn ziface.IConnection){
	fmt.Println("==> DOConnectionLost is Called...")
	fmt.Println("conn ID = ", conn.GetConnID(), " is Lost... " )
}






func main() {
	//1. create server handler, use Zinx Api

	s := znet.NewServer()

	// 2. Register Conn Hook
	s.SetOnConnStart(DOConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 3. Add custmoized Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})


	// customize AddRouter
	//s.AddRouter(&PingRouter{})

	//2. initalize server
	s.Serve()

}
