package net

import (
	"fmt"
	"net"
)

type tcpListen struct {
	address string
	listener *net.TCPListener
	handleConn   func(conn net.Conn)
}

func NewTcpListenManager(address string, handleConn func(conn net.Conn)) *tcpListen {
	return &tcpListen{
		address: address,
		handleConn: handleConn,
	}
}

func (tcpListen *tcpListen) StartTcpListen()  error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", tcpListen.address)
	if err != nil {
		return err
	}
	tcpListen.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("listener", fmt.Errorf("%s", err))
			}
		}()
		for  {
			conn, err := tcpListen.listener.AcceptTCP()
			if err != nil {
				fmt.Println("stop listener", fmt.Errorf("%s", err))
			}
			tcpListen.handleConn(conn)
		}
	}()
	return nil
}

func (tcpListen *tcpListen) close()  {
	tcpListen.listener.Close()
}


