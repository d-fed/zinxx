package znet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"awesomeProject/src/zinx/ziface"
	"awesomeProject/src/zinx/utils"
)

type Connection struct {
	TCPServer ziface.IServer	// which server does this connection belong to
	Conn *net.TCPConn	// connected socket TCP socket
	ConnID uint32 // global unique ConnID/SessionID
	MsgHandler ziface.IMsgHandle	// Message management MsgID and message management module corresponding to processing method
	ctx    context.Context
	cancel context.CancelFunc
	msgChan chan []byte
	msgBuffChan chan []byte

	sync.RWMutex
	property map[string]interface{} // conn property
	propertyLock sync.Mutex // lock to protect current property
	isClosed bool
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	//Initiialize Connection
	c := &Connection{
		TCPServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msgHandler,
		msgChan:     make(chan []byte),
		msgBuffChan: make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:    nil,
	}

	c.TCPServer.GetConnMgr().Add(c) // add Conn to ConnMgr
	return c
}

//StartWriter  Write Data into Client
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				return
			}
			//fmt.Printf("Send data succ! data = %+v\n", data)
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

//StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Reader exit!]")
	defer c.Stop()

	// 创建拆包解包的对象
	for {
		select {
		case <-c.ctx.Done():
			return
		default:

			//读取客户端的Msg head
			headData := make([]byte, c.TCPServer.Packet().GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				fmt.Println("read msg head error ", err)
				return
			}
			//fmt.Printf("read headData %+v\n", headData)

			//拆包，得到msgID 和 datalen 放在msg中
			msg, err := c.TCPServer.Packet().Unpack(headData)
			if err != nil {
				fmt.Println("unpack error ", err)
				return
			}

			//根据 dataLen 读取 data，放在msg.Data中
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					fmt.Println("read msg data error ", err)
					return
				}
			}
			msg.SetData(data)

			//得到当前客户端请求的Request数据
			req := Request{
				conn: c,
				msg:  msg,
			}

			if utils.GlobalObject.WorkerPoolSize > 0 {
				//已经启动工作池机制，将消息交给Worker处理
				c.MsgHandler.SendMsgToTaskQueue(&req)
			} else {
				//从绑定好的消息和对应的处理方法中执行对应的Handle方法
				go c.MsgHandler.DoMsgHandler(&req)
			}
		}
	}
}

//Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TCPServer.CallOnConnStart(c)
}

//Stop 停止连接，结束当前连接状态M
func (c *Connection) Stop() {

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TCPServer.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()
	if c.isClosed == true {
		return
	}
	fmt.Println("Conn Stop()...ConnID = ", c.ConnID) // current Conn is closed
	c.Conn.Close() // Close Socket
	c.cancel()	// Close Writer
	c.TCPServer.GetConnMgr().Remove(c) 	// remove from ConnMgr()
	close(c.msgBuffChan)	// Close All Channel for this Conn
	c.isClosed = true 	// set the mark

}

//GetTCPConnection: obtain the original socket
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *Connection) SendMsg(msgID uint32, data []byte) error { // send to remote TCP
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	dp := c.TCPServer.Packet()		// wrap data and send
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}
	c.msgChan <- msg	// write back to client-side
	return nil
}

func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	dp := c.TCPServer.Packet() 	// packet data and send
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg ID = ", msgID)
		return errors.New("Pack error msg ")
	}
	c.msgBuffChan <- msg	// re-write to client-side
	return nil
}
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// return ctx, and for user to customized go func to obtain
func (c *Connection) Context() context.Context {
	return c.ctx
}



/*
Each incoming request is handled in its own goroutine.
Request handlers often start additional goroutines to access backends such as databases and RPC services.
The set of goroutines working on a request typically needs access to request-specific values
-> identity
-> authorizaiton token
-> request's deadline

When a request is canceled or times out, all the goroutines working on
that request should exit quickly so the system can reclaim any resources they are using.
 */