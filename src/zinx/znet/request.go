package znet

import "awesomeProject/src/zinx/ziface"

// request 包含connection和data
// conn, data
// 把data替换成msg数据结构
type Request struct {
	// current Conn connected with client
	conn ziface.IConnection
	// data requested from client-side
	//data []byte

// 	Add Messge into Request
	msg ziface.IMessage
}

// obtain current conn
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// obtain request's data
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}


func(r* Request) GetMsgID() uint32{
	return r.msg.GetMsgId()
}