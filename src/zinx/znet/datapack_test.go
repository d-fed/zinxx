package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// go func inside handle client-side request
// In charge of datapack pack & unpack unittest
/*
	Stimulate Server
	1. Create socketTCP
	2. Read Data from client-side and do unpack
*/
func TestDataPack(t *testing.T) {
	/*
		simulated server
	*/
	// 1. create socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err :", err)
		return
	}

	/*
		handle request from client-side
	*/

	// 2. create go func to in charge of handling tasks from client-side
	//for{
	// headData := make([]byte, dp.GetHeadLen())
	// dp.Unpack(headData)
	//}
	go func() {
		// read data from client, unpack
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
			}
			// handle client request
			// dp.GetHeadLen()-> headData 读头
			// dp.GetMsgLen() -> msgHead  判断头的长度大于0
			// 1. read from conn read head from pack
			// 2. read from conn, read data from head datalen 【read from client, and unpack】
			go func(conn net.Conn) {
				// --> unpack process <--
				/*
				 define a unpack object dp
				*/
				dp := NewDataPack()
				for {
					// 1. read head from the package from conn
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						break
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						return
						// msgHead belongs to Message struct
						// Id		uint32
						// DataLen	uint32
						// Data		[]byte
					}


					// msg有数据的，第二次读取
					if msgHead.GetMsgLen() > 0 {

						// 2.  based on datalen to read data from conn
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen()) // 切片
						// read from
						// Data	 	[]byte



						// read from io.stream, based on dataLen
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}



						// a completed  message has been read
						fmt.Println("--> Recv MsgID: ", msg.Id, ", datalen =", msg.DataLen, "data = ", msg.Data)

					}
					// 2. read from conn, based on datalen to read data

				}
			}(conn)
		}
	}()
	// 2. obtain data from client, unpack





	/*
		simulate client
	*/
	//客户端goroutine，负责模拟粘包的数据，然后进行发送
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:7777")
		if err != nil {
			fmt.Println("client dial err: ", err)
			return
		}

		/*
			create a pack object dp
		*/
		dp := NewDataPack()

		// simulate package-sticky process, pack 2 msg and send together
		// pack msg1
		msg1 := &Message{
			Id:      1,
			DataLen: 4,
			Data:    []byte{'z', 'i', 'n', 'x'}, // 四个包
		}
		sendData1, err1 := dp.Pack(msg1) // pack
		if err1 != nil {
			fmt.Println("client pack msg1 error", err1)
			return
		}



		// pack msg2
		msg2 := &Message{
			Id:      1,
			DataLen: 5,
			Data:    []byte{'h', 'k', 'n', 'z', 'f', 'o', 'l'}, // 7个包
		}
		sendData2, err2 := dp.Pack(msg2)
		if err2 != nil {
			fmt.Println("client pack msg2 error", err2)
			return
		}

		// append msg1 and msg2
		sendData1 = append(sendData1, sendData2...)
		// parameter ... make variadic parameter
		// accept zero or more string arguments

		// send to server-side
		conn.Write(sendData1)
	}()


	// Client-side Block
	select {
	case <-time.After(time.Second):
		return
	}

	/*
		The ... parameter syntax makes a variadic parameter.
		It will accept zero or more string arguments, and reference them as a slice.
	*/

	// send together to client for msg1 and msg2
}
