package main

//import "zinx-connect/src/zinx/znet"

import "awesomeProject/src/zinx/znet"

//import "zinx-connect/src/zinx/znet"

/*
	create server-side application
*/

func main() {
	//1. create server handler, use Zinx Api
	s := znet.NewServer("[zinxv0.2..]")

	//2. initalize server
	s.Serve()

}
