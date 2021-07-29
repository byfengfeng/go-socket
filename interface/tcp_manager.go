package _interface

type ITcpManager interface {

	StartTcpListen()  error

	Close()
}
