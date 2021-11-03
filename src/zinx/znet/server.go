package znet

// replace Router property into msgHandler Property
import (
	"awesomeProject/src/zinx/utils"
	"awesomeProject/src/zinx/ziface"
	//"errors"
	"fmt"
	"net"
)

// 在server listen之前添加
// 在处理完拆包

type Server struct {
	Name string
	IPVersion string // tcp or others
	IP string
	Port int
	msgHandler ziface.IMsgHandle 	// current server Message management Module, bind MsgID and corresponding Tasks API
	//Router ziface.IRouter	// add router to current Server		// change Router into Message Handler
	ConnMgr ziface.IConnManager					// ConnMgr in current Server
	OnConnStart func(conn ziface.IConnection)	// Hook to start conn
	OnConnStop	func(conn ziface.IConnection)	// Hook to stop conn
}

// define current client connected handleAPI (this handle is fixed TODO: will let client to define handleAPI)
// conn, data from client, cnt: length of package
//func CallBackToClient(conn *net.TCPConn,data []byte, cnt int) error{
//	// echo task
//	fmt.Println("[Conn Handle] CallbackToClient ...")
//	if _,err := conn.Write(data[:cnt]); err !=nil {
//		fmt.Println("write back buf err", err)
//		return errors.New("CallB ackToClient error")
//	}
//
//	return nil
//}

func (s *Server) Start() {
	//fmt.Printf("[Zinx] Server Name : %s, listener at IP: %s, Port:%d is starting", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	//fmt.Print("[Zinx] Version %是， MaxConn: %d, MaxPackeetSize:%d\n",
	//	utils.GlobalObject.Version,
	//	utils.GlobalObject.MaxConn,
	//	utils.GlobalObject.MaxPacketSize,
	//)

	//fmt.Printf("[Start] Server Listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port,
	//	utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	//fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPackeetSize:%d\n",
	//	utils.GlobalObject.Version,
	//	utils.GlobalObject.MaxConn,
	//	utils.GlobalObject.MaxPacketSize)
	// read from Client-side


	fmt.Printf("[START] Server name: %s,listener at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	go func() {
		// 0. start Message Queue, Worker Pool
		s.msgHandler.StartWorkerPool()
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))			// 1. obtain 1 TCP Addr
		if err != nil {
			fmt.Println("resolve tcp address error: ", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr)			// 2. listener to monitor address
		if err != nil {
			panic(err)
		}
		fmt.Println("start Zinx server succ, ", s.Name, " succ, Listenning... ") // successfully listenning to address

		//TODO: generate ConnID automatically
		var cid uint32
		cid = 0

		// 3. block connection from waited client, and handle the client-connection tasks (Read & Write)
		// StartServer
		for {
			conn, err := listener.AcceptTCP() 	// 3.1 AcceptTCP get remote addr
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())


			// 3,2 Set server MaxConn, if execeed, close and start new conn
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++

			// start current conn task
			go dealConn.Start()

			// clients have established connections, do max 512 echo
			// ↓  will be handle by connection module
			go func() {
				// 0. Start Worker Pool
				s.msgHandler.StartWorkerPool()
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}

					// 设置最大连接个数的按断，如果超过个数
					if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
						//TODO 给客户端相应的一个错误超出最大链接的错误包
						conn.Close()
						continue
					}
					fmt.Printf("r ecv client buf %s, cnt %d\n", buf, cnt)
					// connMger加入Server模块，给Server添加一个ConnMgr属性
					// 添加NewServer方法，加入ConnMgr初始化

					// 3.3 handle new conn,，此时应该有handler和Conn绑定
					dealConn := NewConnection(s, conn, cid, s.msgHandler)
					cid++

					// 3.4 启动当前的链接业务处理
					go dealConn.Start()

					// echo data
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}
	}()

}

//Stop 停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	//
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	// monitor & tasks
	s.Start()
	//TODO: extra work

	// Blocking status
	select {}
}

// 添加路由方法
// 给当前服务注册一个路由业务方法，供客户端连接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

/*
	Initialize Server Module Methods / Handler
*/
func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()

	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

