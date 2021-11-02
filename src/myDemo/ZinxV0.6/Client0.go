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
func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)
	// 1. connect remote server, generate conn
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println(  "client start err, exit!", err)
		return
	}

	for {
		// 发送封包Message消息
		dp := znet.NewDataPack()

		// NewMsgPackage
		// convert message to TVL format
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0,[]byte("ZinxV0.8 client Test Message 0000")))
		if err != nil{
			fmt.Println("Pack error:", err)
			return
		}


		// err 的作用域位于if 中，不在上一个TLV的binaryMsg中
		if _, err := conn.Write(binaryMsg); err != nil{
			fmt.Println("write error:", err)
			return
		}

		// Server need to return Message; MsgID=1: ping  ping  ping

		// 1. Read head from stream to retrieve ID and dataLen first
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil{
			fmt.Println("read head error", err)
			break
		}




		// unpacak binaryHead to msg struct(which is msgHead only contains ID and dataLen)
		// 2. binaryHead  --> msgHead
 		//  再根据DataLen 进行第二次读取，将data读出来 Large msg data
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil{
			fmt.Println("client unpack msgHead error", err)
			break
		}


		if msgHead.GetMsgLen() > 0{
			msg := msgHead.(*znet.Message) // type_conversion msgHead --> msg
			msg.Data = make([]byte, msg.GetMsgLen()) 	//2. read data from DataLen

			if _, err := io.ReadFull(conn, msg.Data); err != nil{
				fmt.Println("read msg data error, ", err)
				return
			}

			fmt.Println("---> Recv Server Msg : ID = ", msg.Id, ", len =", msg.DataLen, " , data = ", string(msg.Data))

		}
		// cpu 阻塞
		// 2. conn use Write, write data
		//_, err := conn.Write([]byte("Hello Zinx"))
		//if err != nil {
		//	fmt.Println("write conn err", err)
		//	return
		//}
		//
		//buf := make([]byte, 512)
		//cnt, err := conn.Read(buf)
		//if err != nil {
		//	fmt.Println("read buf error")
		//	return
		//}
		//
		//// length, buf, cnt
		//fmt.Printf(" server call back: %s, cnt = %d\n", buf, cnt)

		// cpu blocking  every 1 second
		time.Sleep(1 * time.Second)

	}

}
