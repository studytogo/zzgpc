package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	port := ":8888"
	startClient(port)
}

func startClient(port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(err)
	}

	//发送信息
	go sendMessage(conn)
	//接收信息
	mes := make([]byte, 1024)
	for {
		lens, err := conn.Read(mes)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			os.Exit(0)
		}
		// if lens > 0 {
		// 	mes[lens] = 0
		// }
		message := string(mes[0:lens])
		fmt.Println(message)
	}
}

func sendMessage(conn net.Conn) {
	var input string
	userName := conn.RemoteAddr().String()
	for {
		fmt.Print("请输入：")
		fmt.Scanln(&input)
		if input == "/quit" {
			conn.Close()
			fmt.Println("BYB...")
			os.Exit(0)
		}
		_, err := conn.Write([]byte(userName + "发送了：：" + input))
		if err != nil {
			fmt.Println(err)
			conn.Close()
			os.Exit(0)
		}
	}
}
