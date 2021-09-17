package nets

import (
	"fmt"
	_interface "game_frame/interface"
	"net"
)

type tcpListen struct {
	address string
	listener *net.TCPListener
	handleConn   func(conn net.Conn)
}

func NewTcpListenManager(address string, handleConn func(conn net.Conn)) _interface.ITcpManager {
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
	defer tcpListen.Close()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("listener", fmt.Errorf("%s", err))
			}
		}()
		for {
			conn, err := tcpListen.listener.AcceptTCP()
			if err != nil {
				fmt.Println("stop listener", fmt.Errorf("%s", err))
			}
			tcpListen.handleConn(conn)
		}
	}()
	fmt.Println("tcp listener start")
	return nil
}

func (tcpListen *tcpListen) Close()  {
	tcpListen.listener.Close()
}


