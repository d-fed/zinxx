package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgID uint32, router IRouter)  // register a router for current service to enable client to handle tasks in connection
	GetConnMgr() IConnManager

	SetOnConnStart(func(IConnection)) 		//Set Server connect Hook
	SetOnConnStop(func(IConnection)) 		//Set server disconnect Hook
	CallOnConnStart(conn IConnection) 		//call Conn OnConnStart Hook
	CallOnConnStop(conn IConnection) 		//call Conn OnConnStop Hook
	Packet() Packet
}

