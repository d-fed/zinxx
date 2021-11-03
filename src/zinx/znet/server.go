package znet

// replace Router property into msgHandler Property
import (
	"awesomeProject/src/zinx/utils"
	"awesomeProject/src/zinx/ziface"
	//"errors"
	"fmt"
	"net"
)

// 在server listen
type Server struct {
	Name       string
	IPVersion  string // tcp or others
	IP         string
	Port       int
	msgHandler ziface.IMsgHandle // current server Message management Module, bind MsgID and corresponding Tasks API
	//Router ziface.IRouter	// add router to current Server		// change Router into Message Handler
	ConnMgr     ziface.IConnManager           // ConnMgr in current Server
	OnConnStart func(conn ziface.IConnection) // Hook to start conn
	OnConnStop  func(conn ziface.IConnection) // Hook to stop conn

	packet ziface.Packet
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
	fmt.Printf("[START] Server name: %s,listener at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)
	go func() {
		s.msgHandler.StartWorkerPool()                                                   // 0. start Message Queue, Worker Pool
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port)) // 1. obtain TCP Addr
		if err != nil {
			fmt.Println("resolve tcp address error: ", err)
			return
		}
		listener, err := net.ListenTCP(s.IPVersion, addr) // 2. monitor address
		if err != nil {
			panic(err)
		}
		fmt.Println("start Zinx server succ, ", s.Name, " succ, Listenning... ")

		//TODO: generate ConnID automatically
		var cid uint32
		cid = 0

		// 3. Start Server COnn
		for {
			conn, err := listener.AcceptTCP() // 3.1 block&wait client to establish
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			// 3,2 Set server MaxConn, if execeed MaxConn, close and start new conn
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			// 3.3 handle new conn request, handler^conn
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++

			// 3.4 go dealConn.Start()
			go dealConn.Start()
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

// AddRouter: server regiester a Router method for client-side
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!")
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// Hook Start & Stop
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("--->CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

func (s *Server) Packet() ziface.Packet{
	return s.packet
}


//func (s *Server) Packet() ziface.Packet{
//	return s.packet
//}

/*
	Initialize Server Module Methods / Handler
*/
//func NewServer(opts ...Option) ziface.IServer {
func NewServer() ziface.IServer {
	utils.GlobalObject.Reload()

	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		packet:     NewDataPack(),
	}
	//for _, opt := range opts {
	//	opt(s)
	//}
		return s
}
