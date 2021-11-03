package main

import (
	"awesomeProject/src/zinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

/*
	stimulate client
*/

/*

10W client请求的数量。
30W请求，
20W 阻塞，并不会浪费资源
- 10W reader
- 10W writer




10 (fixed number) Gouroutine to handle tasks
-> No matter how big client-request
-> There only need 10 goroutine switches

10W处理客户端业务
*/
func main() {
	//time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999") // 1. connect remote server, generate conn
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	for {
		/*
			send packed message
			MsgID = 0
		*/
		dp := znet.NewDataPack() // 发封包Messge的消息
		// convert message to TVL format
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx client Demo Test MsgID=0, [Ping]")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}

		// 先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) // ReadFull 会吧Msg充填满为止
		if err != nil{
			fmt.Println("read head error")
			break
		}

		// 将headDate字节流，拆包到msg中
		// headData => msgHead
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err: ", err)
			return
		}

		if msgHead.GetDataLen()>0{
			// msg是有data数据的，需要再次读取data数据
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// 根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msgData)
			if err != nil{
				fmt.Print("server unpack data err:", err)
				return
			}
			fmt.Println("==> Test Router:[Ping] Recv Msg: =", msgID, ", len= ", msg.DataLen, ", data=")
		}
	}
}
//
//
//		if _, err := conn.Write(binaryMsg); err != nil { // err 的作用域位于if 中，不在上一个TLV的binaryMsg中
//			fmt.Println("write error:", err)
//			return
//		}
//		binaryHead := make([]byte, dp.GetHeadLen())              // 1. Read head from stream
//		if _, err := io.ReadFull(conn, binaryHead); err != nil { //ReadFull, fill msg
//			fmt.Println("read head error", err)
//			break
//		}
//		msgHead, err := dp.Unpack(binaryHead) // put binary head unpack into msg struct=
//		if err != nil {
//			//fmt.Println("client unpack msgHead error", err)
//			fmt.Println("server unpack msgHead error", err)
//			return
//		}
//
//		// msg具有数据
//		if msgHead.GetMsgLen() > 0 { // msg contains data indeed
//			//2. read from DataLen
//			// 根据DataLen做第二次数据读取
//			msg := msgHead.(*znet.Message) // Type Conversion msgHead --> msg
//			msg.Data = make([]byte, msg.GetMsgLen())
//
//			if _, err := io.ReadFull(conn, msg.Data); err != nil {
//				fmt.Println("read msg data error, ", err)
//				return
//			}
//
//			fmt.Println("---> Recv Server Msg : ID = ", msg.Id, ", len =", msg.DataLen, " , data = ", string(msg.Data))
//
//		}
//		// cpu 阻塞
//		// 2. conn use Write, write data
//		//_, err := conn.Write([]byte("Hello Zinx"))
//		//if err != nil {
//		//	fmt.Println("write conn err", err)
//		//	return
//		//}
//		//
//		//buf := make([]byte, 512)
//		//cnt, err := conn.Read(buf)
//		//if err != nil {
//		//	fmt.Println("read buf error")
//		//	return
//		//}
//		//
//		//// length, buf, cnt
//		//fmt.Printf(" server call back: %s, cnt = %d\n", buf, cnt)
//
//		// cpu blocking  every 1 second
//		time.Sleep(1 * time.Second)
//
//	}
//
//}
