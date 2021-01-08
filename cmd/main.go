package main

import (
	"encoding/json"
	"fmt"
	"game_frame/net"
	"time"
)

func main() {
	//var a  = 2048
	////for {
	//	c := string(a)
	//	fmt.Println(len(c))
	//for i:=0; i < len(c); i++ {
	//	fmt.Println(c[i])
	//}
	//	bytes := []byte(c)
	//	if len(bytes) == 3{
	//		for i:=0; i < len(bytes); i++ {
	//			fmt.Println(bytes[i])
	//		}
	//		fmt.Println(len(bytes),"----",a)
	//		//break
	//	}
		//a++
	//}
	//
	go client()
	tcp := net.NewTcpListen("192.168.3.156:30000")
	defer tcp.Close()
	net.SocketListenEvent(tcp, net.Handle)

}

func client()  {
	time.Sleep(time.Second *5)
	tcpDial := net.NewTcpDial("192.168.3.156:30000")
	users := "1"
	for j:=2; j < 50; j++{
		users+=fmt.Sprint(j)
	}
	for i:=1; i <= 10; i++ {
		marshal, _ := json.Marshal(users)
		net.Write(int32(1000+i),marshal,tcpDial)

	}
}
