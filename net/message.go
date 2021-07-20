package net

import "net"

type Message struct {
	Code uint16
	Body []byte
	Connect net.Conn
}
