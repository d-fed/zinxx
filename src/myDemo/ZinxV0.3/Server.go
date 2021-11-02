package main

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

// Test PreRouter
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("Call back before ping error")
	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle ... ")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("callback ping...ping...ping error ")
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle ... ")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping"))
	if err != nil {
		fmt.Println("callback after ping error ")
	}
}

func main() {
	//1. create server handler, use Zinx Api
	s := znet.NewServer("[zinxv0.3..]")

	// customize AddRouter
	s.AddRouter(&PingRouter{})

	//2. initalize server
	s.Serve()

}
