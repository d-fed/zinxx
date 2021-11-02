package znet

import "awesomeProject/src/zinx/ziface"

// implement router, inject BaseRouter, rewrite this base method
type BaseRouter struct {
}

//  Hook before conn
func (br *BaseRouter) PreHandle(request ziface.IRequest) {

}

//  Hook during conn
func (br *BaseRouter) Handle(request ziface.IRequest) {

}

// Hook after conn
func (br *BaseRouter) PostHandle(request ziface.IRequest) {

}
