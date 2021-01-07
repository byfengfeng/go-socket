package net

import (
	"encoding/json"
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
		err := WriteTCP(1,[]byte(strData),conn)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Send:", strData)
	}

	stopChan <- true
}


func WriteTCP(operation byte, w []byte, conn net.Conn) (err error) {
	var ( //operation 包头信息，w包主体信息，
		length = len(string(w))
		buf    = make([]byte, length+5)
	)

	buf[0] = operation
	buf[1] = byte(uint32(length))
	buf[2] = byte(uint32(length) >> 8)
	buf[3] = byte(uint32(length) >> 16)
	buf[4] = byte(uint32(length) >> 24)
	copy(buf[5:], w)
	println("len buf: ", len(buf))
	if _, err = conn.Write(buf); err != nil {
		print("conn.Write(%s), failed(%s)", string(buf), err)
		return
	}
	return
}

type Users struct {
	UserName string
	Age int32
	Sex string
}

func HandleClient(conn net.Conn) {
	defer conn.Close()
	for {
		var (
			buf [502]byte
		)
		_, err := conn.Read(buf[0:5])
		if err != nil {
			return
		}
		lens := int(uint32(buf[1]) | uint32(buf[2])<<8 | uint32(buf[3])<<16 | uint32(buf[4])<<24)
		fmt.Println("lens:", lens)
		from := conn.RemoteAddr()
		fmt.Println("from: ", from)
		switch buf[0] {
		case 1:
			_, err = conn.Read(buf[0:lens])
			if err != nil {
				print("conn.Read() failed(%s)", err)
				continue
			}
			fmt.Println("buf:  ", string(buf[0:lens]))
			user := &Users{}
			json.Unmarshal(buf[0:lens],&user)
			fmt.Println("userName:  ",user.UserName )
			fmt.Println("age:  ",user.Age )
			fmt.Println("sex:  ",user.Sex )

		default:
			print("no data!!")
		}
	}
	return
}
