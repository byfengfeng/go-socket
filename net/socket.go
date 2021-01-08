package net

import (
	"fmt"
	"net"
	"strings"
)




func NewTcpListen(address string)  net.Listener {
	tcp, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("tcp listen failed, err:", err)
		return nil
	}
	return tcp
}

func NewTcpDial(address string)  net.Conn {
	tcp, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("tcp Dial failed, err:", err)
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
func SocketListenEvent(socket net.Listener,channeHandel func(conn net.Conn))  {
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
func ChanneHandel(conn net.Conn) {
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
		err := Write(2301,[]byte(strData),conn)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Send:", strData)
	}

	stopChan <- true
}


func Write(msgCode int32, msgData []byte, conn net.Conn) (err error) {
	var ( //operation 包头信息，msgData包主体信息，
		length = len(string(msgData))
		buf    = make([]byte, length+8)
		code = []byte(fmt.Sprintf("%d",msgCode))
	)

	//验证协议是否合法
	if msgCode < int32(MinCode) || msgCode > int32(MaxCode) {
		return
	}

	//封装协议
	Encode(buf,code,msgData,length)

	//发送消息
	if _, err = conn.Write(buf); err != nil {
		print("conn.Write(%s), failed(%s)", string(buf), err)
		return
	}
	return
}

func Handle(conn net.Conn) {
	defer conn.Close()
	for {
		//解析消息号
		msgCode, lens, err := DecodeRead(conn)
		if err != nil {
			return
		}
		var buf = make([]byte,lens)
		_, err = conn.Read(buf[0:lens])
		if err != nil {
			fmt.Printf("conn.Read() failed(%s)", err)
			continue
		}

		//将消息进行封装转发至通道
		msg := Message{
			Code: msgCode,
			Body: buf[0:lens],
		}
		fmt.Println(msg.Code, string(msg.Body))
	}
	return
}
