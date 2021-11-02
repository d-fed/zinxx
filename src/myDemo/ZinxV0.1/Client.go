package main

import (
	"fmt"
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
	conn, err := net.Dial("tcp", "127.0.0.1:8998")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 2. conn use Write, write data
		_, err := conn.Write([]byte("Hello Zinx"))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}

		// length, buf, cnt
		fmt.Printf(" server call back: %s, cnt = %d\n", buf, cnt)

		// cpu blocking  every 1 second
		time.Sleep(1 * time.Second)

	}

}
