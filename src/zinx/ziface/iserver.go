package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgID uint32, router IRouter)  // register a router for current service to enable client to handle tasks in connection
	GetConnMgr() IConnManager
}

