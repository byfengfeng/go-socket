package net

import "net"

type Message struct {
	Code int32
	Body []byte
	Connect net.Conn
}
