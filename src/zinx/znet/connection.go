package znet

// 1. StartReader()
// 2. request
// 3. Execute 3 Router Handles
import (
	"awesomeProject/src/zinx/utils"
	"awesomeProject/src/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	TCPServer ziface.IServer
	Conn *net.TCPConn		// current conn socket TCP
	ConnID uint32 // Conn ID
	isClosed bool //	Current Conn Status
	//handleAPI ziface.HandleFunc
	ExitChan chan bool  		// info current conn has exit or paused channel
	//RW between Goroutine
	msgChan  chan []byte // non-buffer Task channel
	msgBuffChan  chan []byte // Buffer Task channel

	// 消息的管理MsgID 和 对应的处理业务API关系
	MsgHandler ziface.IMsgHandle

	//conn property collection
	property map[string]interface{}


	//READWRITE protect conn
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
		property: nil,
		//msgBuffChan:make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
	}

	// add Conn into  ConnMgr
	c.TCPServer.GetConnMgr().Add(c)
	// NewConnection
	return c
}

// conn Reader 	// 1. Read Client Data 2. Call HandleAPI
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, " [Reader is exit!], remote addr is ", c.RemoteAddr().String())
	// if any break triggered
	defer c.Stop() // 资源回收

	// 1. Read Client Data 2. Call HandleAPI
	for {

		// read Client data into buf, max 512 byte
		//buf := make([]byte, utils.GlobalObject.MaxPacketSize) // buf is from client-side
		////cnt, err := c.Conn.Read(buf)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		// Use Router to replace HandleAPI
		// 调用当前链接锁绑定的HandleAPI
		//(客户端Conn, buf, cnt) buf and buf_length
		//if err := c.handleAPI(c.Conn, buf, cnt); err != nil{
		//	fmt.Println("ConnID ", c.ConnID, " handle is error", err)
		//	//break
		//	continue
		//}

		/*
			1. Create a unpack Object
		*/
		dp := NewDataPack()

		// 2. Read Msg HeadData from client 8 bytes
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		// 对于这8个字节的headData进行解包
		// Unpack to obtain msgID and msgDataLen --> load into --> msg
		msg, err := dp.Unpack(headData) // msg contains msgID and dataLen
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		// Read Data based on data --->load into---> msgDatalen
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)
		// and put inside msg.Data
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// The work pool mechanism has been activated, and the message is handed over to the Worker for processing
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req) // Execute the corresponding Handle method in the bound message and the corresponding processing method
		}

		////// Registe Router
		////// router which implement register
		//// ✅ Correct Way ✅
		//go func(request ziface.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)
		////

		// 得到当前conn数据的Request请求数据
		// because buf has been put into Request
		// Request data from current conn
		//req := Request{
		//	conn: c,
		//	data: msg,
		//}

		// find router corresponding to binding Conn
		//c.Router.PreHandle(req) // Wrong!!!, you need to pass address instead of request

		// remove CallBackToClient method

		// obtain the correspoding router of the Conn
	}

}

/*
 	write message Goroutine, Users send Message
	dedicated to send message to client

*/

/*
	dedicated module which send message to client-side
	write Message function
*/

func (c *Connection) StartWriter() {
	fmt.Print("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")

	// 不断阻塞的等待Channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}

		case <-c.ExitChan: //可读
			// 代表Reader已经退出，此时Writer也已经退出

			return
		}
	}

	// 告知当前链接已退出，Writer退出

}

// start conn, let current conn ready to work
func (c *Connection) Start() {
	fmt.Print("Conn Start() ... ConnID = ", c.ConnID)
	// start READ data tasks from current conn
	go c.StartReader()
	go c.StartWriter()
	c.TCPServer.CallOnConnStart(c) // create conn and handle task, execute Hook function


}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop() .. ConnID = ", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.Conn.Close()                     // close socket
	c.ExitChan <- true                 //close Writer Goroutine
	c.TCPServer.GetConnMgr().Remove(c) // remove from current ConnMgr
	//Exit channel and recycle resources
	close(c.ExitChan)
	close(c.msgChan)

}

// current binding socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// current Conn ID of current module
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// TCP status IP port of remote conn module
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// pack the data(going to send to client-side)
// provide a function called SendMsg function to send data to client-side and
//pack and send
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	// pack data and MsgDataLen|MsgID|Data
	dp := NewDataPack()

	/*
		binaryMsg: serialized message
			MsgDataLen|MsgID|Data
	*/
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发给客户端
	c.msgChan <- binaryMsg
	// Writer得到msgChan，然后发布出去
	if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id ", msgId, " error :", err)
		return errors.New("conn Write error")
	}
	return nil
}


func (c *Connection) SetProperty(key string, value interface{}){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// add a conn
	c.property[key] = value
}


func (c *Connection) GetProperty(key string)(interface{}, error){
	c.propertyLock.RLock()
	defer c.propertyLock.RLock()

	// read property
	if value, ok := c.property[key]; ok{
		return value,nil
	}else{
		return nil, errors.New("no property found")
	}
}



func (c *Connection) RemoveProperty(key string)(){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// delete property
	delete(c.property, key)

}