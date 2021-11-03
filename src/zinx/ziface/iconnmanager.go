package ziface


/*
连接管理模块抽象层
 */

type IConnManager  interface{
	// add conn
	Add(connection IConnection)

	// delete conn
	Remove(conn IConnection)

	Get(connID uint32) (IConnection, error)// obtain conn from connID

	Len() int  // num of connection

	ClearConn()


}