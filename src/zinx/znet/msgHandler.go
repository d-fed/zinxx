package znet

import (
	"awesomeProject/src/zinx/utils"
	"awesomeProject/src/zinx/ziface"
	"fmt"
	"strconv"
)

/*
Message Layer Abstraction
*/
type MsgHandle struct {
	// 存放每个MsgID所对应处理方法
	Apis           map[uint32]ziface.IRouter // store msgID, and map property to handle method
	WorkerPoolSize uint32                    // num of workers  全局配置的参数中获取
	TaskQueue      []chan ziface.IRequest    // MQueue in Worker to obtain task
}

// 每一个Worker都对应一个Go来承载
// 阻塞等待和当前Worker对应的Channel

// 初始化、创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取，也可以在配置文件中让用户设置即可
		// worker 和 queue 一一对应
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen), // request, MaxWorkerTaskLen
		// request, MaxWorkerTaskLen
		// WorkerPoolSize

		// 创建Worker工作池

	}
}

/*
 Request包含：
 	conn ziface.IConnection
 // 当前客户端已经发出的请求
	msg ziface.IMessage
*/

// 索引 / 调度 / 执行对应Router消息的处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1. find msgID from Request
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND! Need Register!")
		return
	}

	// 2. call Router based on MsgID
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// request里面有conn，msgID和对应的date
// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1.Determine whether the API processing method bound to the current msg already exists
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api , msgID = " + strconv.Itoa(int(msgID)))
	}
	// 2. msg+api binding
	mh.Apis[msgID] = router
	fmt.Println("Add api msgID = ", msgID)
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started .... ")

	// 不断的阻塞等待对应消息
	for {
		select {
		// taskQueue过来， call DoMsgHandle
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
			// 保证每一个Worker收到的Request任务是均衡的
		}
	}

	// 每一个Worker利用Go Func
}

// 将消息平均分配给不同的Worker
// 将消息交给TaskQueue，由Worker进行处理
// 根据IP的地域进行分配，按照轮询分配
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1. 将消息平均分配给不同的Worker  	// 根据客户端建立的ConnID来分配  	// 基本的平均分配的轮询法则

	// 得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	// 当前需要工作的ID数量
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), " request MsgID = ",
		request.GetMsgID(),
		" to WorkerID = ", workerID)
	// 将请求消息发给任务队列
	mh.TaskQueue[workerID] <- request

	// 将消息队列Queue集成到Zinx框架中
	// 2. 将消息发送给对应的Worker的TaskQueue

	// 对线程池的数量取余

	// 2. 将消息发送给对应Worker的TaskQueue
}

// 1. Initialzie Worker Pool
// 一个Zinx 框架只能有一个工作池
// 2. Start Worker working Process, block and wait, MessageTask
// 对外暴露：func (mh* MsgHandle) StartOneWorker(){
func (mh *MsgHandle) StartWorkerPool() {
	// 根据WorkerPoolSize 分别开启Worker，每一个Worker用go 来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个Worker启动，每一个Worker对应一个Channel，开辟管道空间
		// Worker ID 和 Channel ID对应
		// 1. 当前的Worker对应Channel消息队列，开辟空间，第0个worker，就用第0个Channel

		// 每一个Worker关联一个TaskQueue
		// WorkerPoolSize
		// MaxWorkerTaskLen=管道长度，每一个Worker对应的消息队列的任务的数量的最大值
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen) // MaxWorkerTaskLen 请求的最大长度

		// 2. 启动当前的Worker，阻塞等待消息从Channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
		// 当前队列的编号，和mh.TaskQueue[i]
	}
}
