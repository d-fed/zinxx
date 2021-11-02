package ziface

type IServer interface {
	// 启动服务器
	Start()

	// 停止服务器
	Stop()

	// 运行服务器
	Serve()

	// register a router  for current service to enable client to handle the connection
	AddRouter(msgID uint32, router IRouter)
}
