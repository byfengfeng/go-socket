package tcp_channel

import (
	"errors"
	"net"
)

type channel struct {
	conn net.Conn
	readChan chan []byte
	writeChan chan []byte
}

func (c *channel) write(bytes []byte) error {
	_, err := c.conn.Write(bytes)
	return err
}

func (c *channel) read() (err error) {
	var(
		headByte []byte
		length uint16
	)
	for  {
		headByte = make([]byte,2)
		_, err = c.conn.Read(headByte)
		if err != nil {
			return err
		}
		length = uint16(headByte[0]) << 8 | uint16(headByte[1])

		if length <= 0{
			return errors.New("错误消息长度")
		}
		body := make([]byte,length-2)
		_, err = c.conn.Read(body)
		if err != nil {
			return err
		}
		c.readChan <- body
	}
	if err != nil {
		return
	}
	return nil
}