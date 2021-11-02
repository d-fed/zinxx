package znet

// replace Router property into MsgHandler Property
import (
	"awesomeProject/src/zinx/utils"
	"awesomeProject/src/zinx/ziface"
	//"errors"
	"fmt"
	"net"
)


// 在server listen之前添加
// 在处理完拆包



//iServer 的接口实现，定义一个Server的服务器module
type Server struct {
	//name
	Name string
	//ip version
	IPVersion string
	//IP listen to
	IP string
	//port listen to
	Port int

	// current server Message management Module, bind MsgID and corresponding Tasks API
	MsgHandler ziface.IMsgHandle

	// add router to current Server
	//Router ziface.IRouter
	// change Router into Message Handler

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
	fmt.Printf("[Zinx] Server Name : %s, listenner at IP: %s, Port:%d is starting", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Print("[Zinx] Version %是， MaxConn: %d, MaxPackeetSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize,
	)

	fmt.Printf("[Start] Server Listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port,
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPackeetSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)
	// read from Client-side
	go func() {
		// 0. 开启消息队列，以及Worker工作池
		// 0. start Message Queue, Worker Pool
		s.MsgHandler.StartWorkerPool()


		// 1. obtain 1 TCP Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error: ", err)
			return
		}
		// Success obtain the   addr

		// 2. [obtain listenner from addr] monitor server address
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}

		fmt.Println("start Zinx server succ, ", s.Name, " succ, Listenning... ")
		var cid uint32
		cid = 0

		// 3. block connection from waited client, and handle the client-connection tasks (Read & Write)
		for {
			// return conn once receive from client
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 在Server中建立New Connection
			// bind conn and latest conn tasks to generate conn module
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			// start current conn task
			go dealConn.Start()

			// clients have established connections, do max 512 echo
			// ↓  will be handle by connection module
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}

					fmt.Printf("r ecv client buf %s, cnt %d\n", buf, cnt)

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

func (s *Server) Stop() {
	//TODO: open resources, status， stop and recycle
	// connection info to pause or recycle

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
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!!")
}

/*
	Initialize Server Module Methods
*/

func NewServer(name string) ziface.IServer {
	s := &Server{
		//Name :name,
		//IPVersion: "tcp4",
		//IP: "0.0.0.0",
		//Port:8998,
		//Router: nil,
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
