package main

import (
	"encoding/json"
	"fmt"
	"game_frame/net"
	"time"
)

func main() {
	go asd()
	tcp := net.NewTcpListen("192.168.3.156:30000")
	defer tcp.Close()
	net.SocketListenEvent(tcp, net.HandleClient)

}

func asd()  {
	fmt.Println("111")
	time.Sleep(time.Second *5)
	tcpDial := net.NewTcpDial("192.168.3.156:30000")
	users := net.Users{
		UserName: "张三",
		Age:      123,
		Sex:      "男",
	}
	for i:=1; i <= 10; i++ {
		marshal, _ := json.Marshal(users)
		net.WriteTCP(1,marshal,tcpDial)
	}
}
