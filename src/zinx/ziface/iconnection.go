package ziface

import "net"

// define the abstraction layer of connection
type IConnection interface {
	// start connection, start work
	Start()

	// end connection, stop current work
	Stop()

	// obtain current socket conn
	GetTCPConnection() *net.TCPConn

	// obtain current module conn ID
	GetConnID() uint32

	// obtain Client TCP IP + Port
	RemoteAddr() net.Addr

	// send data to remote server
	SendMsg(msgID uint32, data []byte) error 		//直接将Message数据发送数据给远程的TCP客户端(无缓冲)
	// send data to r
}

// define a method to handle the linked task
type HandleFunc func(*net.TCPConn, []byte, int) error
