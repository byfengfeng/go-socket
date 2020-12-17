package net

import (
	"fmt"
	"net"
	"strings"
)

func NewTcp(address string)  net.Listener {
	tcp, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("tcp listen failed, err:", err)
		return nil
	}
	return tcp
}

func NewUdp(address string)  net.Listener {
	udp, err := net.Listen("udp", address)
	if err != nil {
		fmt.Println("udp listen failed, err:", err)
		return nil
	}
	return udp
}

//socket监听事件
func socketListenEvent(socket net.Listener,channeHandel func(conn net.Conn))  {
	defer socket.Close()
	for {
		conn, err := socket.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go channeHandel(conn)
	}
}


//处理每个连接用户的链接
func channeHandel(conn net.Conn) {
	defer conn.Close()

	readChan := make(chan string)
	writeChan := make(chan string)
	stopChan := make(chan bool)

	go readConn(conn, readChan, stopChan)
	go writeConn(conn, writeChan, stopChan)

	for {
		select {
		case readStr := <-readChan:
			upper := strings.ToUpper(readStr)
			writeChan <- upper
		case stop := <-stopChan:
			if stop {
				break
			}
		}
	}
}

//读取
func readConn(conn net.Conn, readChan chan<- string, stopChan chan<- bool) {
	for {
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			break
		}

		strData := string(data)
		fmt.Println("Received:", strData)

		readChan <- strData
	}

	stopChan <- true
}


//写入
func writeConn(conn net.Conn, writeChan <-chan string, stopChan chan<- bool) {
	for {
		strData := <-writeChan
		_, err := conn.Write([]byte(strData))
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Send:", strData)
	}

	stopChan <- true
}
