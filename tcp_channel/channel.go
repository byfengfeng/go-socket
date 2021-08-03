package tcp_channel

import "net"

type channel struct {
	conn net.Conn
	readChan chan []byte
	writeChan chan []byte
}

func (c *channel) write(bytes []byte) error {
	_, err := c.conn.Write(bytes)
	return err
}

func (c *channel) read() error {

	return nil
}