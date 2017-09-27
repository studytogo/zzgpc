package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	port := ":8888"
	startServer(port)
}

func startServer(port string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	if err != nil {
		fmt.Println(err)
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
	}

	//创建接收信息的map
	//channel管道
	conns := make(map[string]net.Conn)
	message := make(chan string, 10)

	//将接受的信息发送给客户端
	go sendMesClient(&conns, message)

	fmt.Println("准备连接客户端。。。。")
	for {
		time.Sleep(time.Second * 5)
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println(conn.RemoteAddr().String(), "连接成功")

		conns[conn.RemoteAddr().String()] = conn
		go serverAcceptMes(conn, message)
	}
}

func serverAcceptMes(conn net.Conn, message chan string) {
	buf := make([]byte, 1024)
	for {
		lens, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			conn.Close()
			os.Exit(0)
		}
		if lens > 0 {
			buf[lens] = 0
		}
		message <- string(buf[0:lens])
	}
}

func sendMesClient(conns *map[string]net.Conn, message chan string) {
	for {
		mes := <-message
		fmt.Println(mes)
		for key, v := range *conns {
			_, err := v.Write([]byte(key + "say::" + mes))
			if err != nil {
				fmt.Println(err.Error())
				delete(*conns, key)
			}
		}
	}
}
