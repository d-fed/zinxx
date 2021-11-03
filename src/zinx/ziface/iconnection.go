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
	SendMsg(msgID uint32, data []byte) error 		//send Msg to remote TCP client 【non-buffer】
}

// define a method to handle the linked task
type HandleFunc func(*net.TCPConn, []byte, int) error
