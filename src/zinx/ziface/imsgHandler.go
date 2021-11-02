package ziface

type IMsgHandle interface {
	DoMsgHandler(request IRequest)          // //调度，执行对应Router消息处理方法（非阻塞）
	AddRouter(msgID uint32, router IRouter) //  为消息添加具体的处理逻辑
	StartWorkerPool() // 启动Worker工作池
	// 根据pool size，对每一个Worker用Go 来承载
	//for i:=0; i < int(mh.WorkerPoolSize); i++{
	SendMsgToTaskQueue(request IRequest) 	// 将消息交给TaskQueue，由Worker进行处理
}
