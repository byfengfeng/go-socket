package channel

import (
	"errors"
	"net"
)

type Channel struct {
	Conn net.Conn
	ConnState ConnState
	ReadChan chan []byte
	WriteChan chan []byte
}

type ConnState uint8

const(
	Online ConnState = iota
	Offline
	Disconnected
)

func NewChannel(conn net.Conn) *Channel {
	return &Channel{
		Conn: conn,
		ReadChan: make(chan []byte),
		WriteChan: make(chan []byte),
		ConnState: Online,
	}
}

func (c *Channel) write() (err error) {
	defer c.Close()
	for  {
		bytes := <- c.WriteChan
		if len(bytes) > 0 {
			_, err = c.Conn.Write(bytes)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Channel) read() (err error) {
	defer c.Close()
	var(
		headByte []byte
		length uint16
	)
	for  {
		headByte = make([]byte,2)
		_, err = c.Conn.Read(headByte)
		if err != nil {
			return err
		}
		length = uint16(headByte[0]) << 8 | uint16(headByte[1])

		if length <= 0{
			return errors.New("错误消息长度")
		}
		body := make([]byte,length-2)
		_, err = c.Conn.Read(body)
		if err != nil {
			return err
		}
		c.ReadChan <- body
	}
}

func (c *Channel) Start()  {
	go c.read()
	go c.write()
}

func (c *Channel) Close()  {
	close(c.ReadChan)
	close(c.WriteChan)
	c.Conn.Close()
}