package znet

import (
	"awesomeProject/src/zinx/ziface"
	"errors"
	"fmt"
	"sync"
)

/*
conn Manager Module
*/
// 链接管理模块
type ConnManager struct {
	connections map[uint32]ziface.IConnection	// 管理的链接集合
	connLock    sync.RWMutex 					// 保护连接集合的读写锁
}


// 创建当前链接的方法
func NewConnManager() *ConnManager{
	return &ConnManager{
		connections: make(map[uint32] ziface.IConnection),
	}
}



// Add 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection){
	connMgr.connLock.Lock() 	// 保护共享资源map, 添加写锁
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn  	// 将conn加入到ConnManager中
	connMgr.connLock.Unlock()
	fmt.Println("connID = ", connMgr.Len())


}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection){
	// 保护共享资源
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除链接信息
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println("connection add to ")
}

// 根据connID获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error){
	connMgr.connLock.RLock()		// 保护共享资源map，加读锁
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}
		return nil, errors.New("connection not FOUND!")

}

// Len得到当前链接的总数
func (connMgr *ConnManager) Len() int{
	connMgr.connLock.RLock()
	length := len(connMgr.connections)
	connMgr.connLock.RUnlock()
	return length
}


//ClearConn 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	connMgr.connLock.Unlock()
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}

//ClearOneConn  obtain connID through connID
func (connMgr *ConnManager) ClearOneConn(connID uint32) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connections := connMgr.connections
	if conn, ok := connections[connID]; ok {
		conn.Stop() 	// stop
		delete(connections, connID) // delete(map, key)
		fmt.Println("Clear Connections ID:  ", connID, "succeed")
		return
	}

	fmt.Println("Clear Connections ID:  ", connID, "err")
	return
}