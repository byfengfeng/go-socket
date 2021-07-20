package net

import (
	"fmt"
	"net"
)

var (
	req chan *Message
	res chan *Message
)

func NewTcpListen(address string,reqMsg chan *Message,resMsg chan *Message)  net.Listener {
	tcp, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("tcp listen failed, err:", err)
		return nil
	}
	req = reqMsg
	res = resMsg
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
		go Write()
	}
}

func Write() (err error) {
	for  {
		m := <-res
		if m != nil {
			var ( //msgCode 包头信息，msgData包主体信息，
				length = uint16(len(string(m.Body)))

				buf    = make([]byte, length+4)
				//bytes  = make([]byte,2)
			)

			//验证协议是否合法
			if m.Code < uint16(MinCode) || m.Code > uint16(MaxCode+2) {
				return
			}
			//bytes = append(bytes[0:],byte(m.Code>>8),byte(m.Code))
			//封装协议
			Encode(m.Code, buf, m.Body, length)
			//发送消息
			if _, err = m.Connect.Write(buf); err != nil {
				print("conn.Write(%s), failed(%s)", string(buf), err)
				return
			}
		}
	}
	//return
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
		req <- &Message{
			Code: msgCode,
			Body: buf[0:lens],
			Connect: conn,
		}
	}

}
