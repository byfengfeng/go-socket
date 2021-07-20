package net

import (
	"net"
)

type MsgCode int32

const (
	MinCode MsgCode = 1000
	MaxCode = 9999
)

func DecodeRead(conn net.Conn) (code uint16,lens uint16,err error) {
	var(
		buf [4]byte
	)
	_, err = conn.Read(buf[0:])
	if err != nil {
		return
	}
	//获取消息码
	code = uint16(buf[0])<<8 | uint16(buf[1])
	//获取数据长度
	lens = uint16(buf[2])<<8 | uint16(buf[3])
	return
}

func Encode(code uint16,buf,msgData []byte,msgLength uint16)  {
	//封装协议号
	buf = append(buf,byte(code>>8),byte(code))
	//封装长度
	buf = append(buf,byte(msgLength>>8),byte(msgLength))

	//深度拷贝消息体
	copy(buf[2:], msgData)
}
