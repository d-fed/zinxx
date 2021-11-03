package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgID uint32, router IRouter)  // register a router for current service to enable client to handle tasks in connection
	GetConnMgr() IConnManager

	SetOnConnStart(func(IConnection)) 		//设置该Server的连接创建时Hook函数
	SetOnConnStop(func(IConnection)) 		//设置该Server的连接断开时的Hook函数
	CallOnConnStart(conn IConnection) 		//调用连接OnConnStart Hook函数
	CallOnConnStop(conn IConnection) 		//调用连接OnConnStop Hook函数
	Packet() Packet
}

