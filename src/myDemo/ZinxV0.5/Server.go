package main

// 只需要保留一个handle即可

//import "zinx-connect/src/zinx/znet"

import (
	"awesomeProject/src/zinx/ziface"
	"awesomeProject/src/zinx/znet"
	"fmt"
)

//import "zinx-connect/src/zinx/znet"

/*
	create server-side application
*/

// ping test customize Router

// inherit base Router
type PingRouter struct {
	znet.BaseRouter
}


// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ... ")


	// read from client-side, and ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(), ", data = ", string(request.GetData()))


	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil{
		fmt.Println(err)
	}
}


func main() {
	//1. create server handler, use Zinx Api
	s := znet.NewServer("[zinxv0.5]")

	// customize AddRouter
	//s.AddRouter(&PingRouter{})

	//2. initalize server
	s.Serve()

}
