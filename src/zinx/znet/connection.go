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
)

type Connection struct {
	// current conn socket TCP
	Conn *net.TCPConn

	// Conn ID
	ConnID uint32

	// Current Conn Status
	isClosed bool

	// Current conn bi nding tasks method
	handleAPI ziface.HandleFunc

	// info current conn has exit or paused channel
	ExitChan chan bool //退出消息

	// 无缓冲id管道，用于读，写Goroutine之间的消息
	msgChan chan [] byte // 业务消息

	// 消息的管理MsgID 和 对应的处理业务API关系
	MsgHandler ziface.IMsgHandle


	// call Router in StartReader
}

//conn, id, router
// initialize conn module
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		MsgHandler:   msgHandler,
		isClosed: false,
		msgChan: make(chan []byte),
		ExitChan: make(chan bool, 1),
	}

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
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize) // buf is from client-side
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
		// 做一个判断，判断是是否有已经开启Worker

		if utils.GlobalObject.WorkerPoolSize > 0{
			// has initialized Worker Pool, send message to Worker Pool to handle
			c.MsgHandler.SendMsgToTaskQueue(&req)
		}else{
			// 从Router中，找到绑定的Conn对应的Router调用
			// 根据绑定好的MsgID，找到对应的api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}



		// 根据绑定好的MsgID，找到对应处理api业务，执行
		go c.MsgHandler.DoMsgHandler(&req)

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


func (c *Connection) StartWriter(){
	fmt.Print("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")


	// 不断阻塞的等待Channel的消息，进行写给客户端
	for{
		select{
		case data := <-c.msgChan:
		// 有数据要写给客户端
		if _,err := c.Conn.Write(data); err != nil{
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

	//TODO: 启动从当前连接的写数据的业务
	go c.StartWriter()

}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop() .. ConnID =", c.ConnID)

	// if close
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// close socket
	c.Conn.Close()

	//告知Writer关闭，退出
	c.ExitChan <- true

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
