package net

import (
	"net"
	"strconv"
)

type MsgCode int32

const (
	MinCode MsgCode = 1000
	MaxCode = 9999
)

func DecodeRead(conn net.Conn) (code int32,lens int,err error) {
	var(
		buf [8]byte
	)
	_, err = conn.Read(buf[0:8])
	if err != nil {
		return
	}
	//获取消息码
	msgCode,_ := strconv.ParseInt(string(buf[0:4]),10,64)
	code = int32(msgCode)
	//获取数据长度
	lens = int(uint32(buf[4]) | uint32(buf[5])<<8 | uint32(buf[6])<<16 | uint32(buf[7])<<24)
	return
}

func Encode(buf,code,msgData []byte,msgLength int)  {
	//封装协议号
	buf[0] = code[0]
	buf[1] = code[1]
	buf[2] = code[2]
	buf[3] = code[3]

	//封装消息体
	buf[4] = byte(uint32(msgLength))
	buf[5] = byte(uint32(msgLength) >> 8)
	buf[6] = byte(uint32(msgLength) >> 16)
	buf[7] = byte(uint32(msgLength) >> 24)

	//深度拷贝消息体
	copy(buf[8:], msgData)
}
