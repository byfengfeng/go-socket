package main

import (
	"fmt"
	"game_frame/net"
	pb "game_frame/probuf"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
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
	//	a++
	//}
	//

	//TCP测试
	//reqChan := make(chan *net.Message)
	//resChan := make(chan *net.Message)
	//go client(resChan)
	//tcp := net.NewTcpListen("192.168.3.156:30000",reqChan,resChan)
	//go func(){
	//	for  {
	//		v := <- reqChan
	//
	//		fmt.Println(v)
	//		notifi := &pb.ResMailChange{
	//		}
	//		proto.Unmarshal(v.Body,notifi)
	//		fmt.Println(notifi.MailId)
	//		fmt.Println(notifi.MailChangeType)
	//		//v.Code = v.Code + 10
	//		//v.Body = []byte(fmt.Sprint("回复消息"))
	//		//resChan <- v
	//	}
	//
	//}()
	//defer tcp.Close()
	//net.SocketListenEvent(tcp, net.Handle)

	//WEBSOCKET测试
	http.Handle("/upper", websocket.Handler(upper))
	http.HandleFunc("/", index)

	if err := http.ListenAndServe(":9999", nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var i = 0
func upper(ws *websocket.Conn) {
	var err error
	i++
	var cnt = i
	connectMap[fmt.Sprint(i)] = ws

	for {
		var reply string
		if len(connectMap) > 1{
			websocket.Message.Send(connectMap["1"], strings.ToUpper(fmt.Sprint(len(connectMap))))
		}
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println(err)
			if err == io.EOF{
				fmt.Println(string(cnt),"下线了")
				delete(connectMap,fmt.Sprint(cnt-1))
			}
			break
		}

		if err = websocket.Message.Send(ws, strings.ToUpper(reply+"123")); err != nil {
			fmt.Println(err)
		}else {
			websocket.Message.Send(ws, strings.ToUpper("12312312"))
		}

	}
}

var connectMap = make(map[string]*websocket.Conn )


func index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func getbys(data interface{}) ([]byte, error) {
	return proto.Marshal(data.(proto.Message))
}

//func main() {
	//code := 13001
	//notifi := &pb.ResMailChange{
	//	MailChangeType: int32(3),
	//	MailId: []int64{23001},
	//}
	////bytes, err := getbys(notifi)
	////if err != nil{
	////	panic(err)
	////}
	////size := notifi.XXX_Size()
	////a, err := notifi.Marshal()
	//marshal, _ := proto.Marshal(notifi)
	//bytes, _ := json.Marshal(notifi)
	////notifi.ProtoMessage()
	////bytes, ints := notifi.Descriptor()
	//fmt.Println(len(marshal))
	//fmt.Println(len(bytes))
	//fmt.Println(len(bytes))

	//http.Handle("/upper", websocket.Handler(upper))
	//http.HandleFunc("/", index)
	//
	//if err := http.ListenAndServe(":9999", nil); err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	//router := mux.NewRouter()
	//router.HandleFunc("/ws", Myws)
	//if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
	//	fmt.Println("err:", err)
	//}
	//go h.run()

//}



func client(data chan *net.Message)  {
	time.Sleep(time.Second *3)
	tcpDial := net.NewTcpDial("192.168.3.156:30000")
	notifi := &pb.ResMailChange{
			MailChangeType: int32(3),
			MailId: []int64{23001},
		}
	marshal, _ := proto.Marshal(notifi)
	//users := "1"
	//for j:=2; j < 50; j++{
	//	users+=fmt.Sprint(j)
	//}
	for i:=1; i <= 10; i++ {
		//marshal, _ := json.Marshal(users)
		data <- &net.Message{
			Code: int32(1000+i),
			Connect: tcpDial,
			Body: marshal,
		}
	}
	for  {
		net.Handle(tcpDial)
	}
}
