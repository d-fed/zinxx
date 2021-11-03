package utils

// Unmarshal: load json into struct
import (
	"awesomeProject/src/zinx/ziface"
	"encoding/json"
	"io/ioutil"
)

/*
存储一切有关Zinx框架的全局参数，供其他模块使用
一些参数可以通过zinx.json由用
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer // current Zinx global Server object
	Host      string         // current server listening IP
	TcpPort   int            // current server listening port
	Name      string         // current server name

	/*
		Zinx
	*/
	Version        string 
	MaxConn        int    // Max allowed server connection
	MaxPacketSize uint32 // Max value of Zinx data packet
	WorkerPoolSize uint32	// current worker pool number of tasks
	// 当前业务工作Worker Pool的 Goroutine 的数量
	MaxWorkerTaskLen uint32	// The max task storage  of the corresponding task responsibilities of the staff
	// MaxWorkerTaskLen，允许用户最多开辟多少个Worker（限定条件下）
	//MaxMsgChanLen
	/*
		config file path
	*/
	ConfFilePath string
}

/*
	Define global object
*/
var GlobalObject *GlobalObj

/*

 */
func init() {
	// if config didn't load, here is the default value
	GlobalObject = &GlobalObj{
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPacketSize: 4096,
	}

	//
}

/*
 from zinx.json to load customized parameter
*/
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}

	// json file to struct
	err = json.Unmarshal(data, &GlobalObject)
}

/*
	Provide init method, initialize current GlobalObject
*/
func init() {
	GlobalObject = &GlobalObj{
		// if config not loaded
		Name:           "ZinxServerApp",
		Version:        "V0.7",
		TcpPort:        8999,
		Host:           "0.0.0.0",
		MaxConn:        1000,
		MaxPacketSize: 4096,
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024, // max len
	}

	// try to use conf/zinx.json to load user's defined
	GlobalObject.Reload()

}
