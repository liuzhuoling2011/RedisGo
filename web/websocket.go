package web

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"golang.org/x/net/websocket"
)

func WSHandler(conn *websocket.Conn) {
	fmt.Printf("Websocket新建连接: %s -> %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())

	for {
		var reply string
		if err := websocket.Message.Receive(conn, &reply); err != nil {
			fmt.Println("Websocket连接断开:", err.Error())
			conn.Close()
			return
		}
		rJson, err := simplejson.NewJson([]byte(reply))
		if err != nil {
			fmt.Println("receive err:", err.Error())
			return
		}
		rType, _ := rJson.Get("type").Int()
		fmt.Print(rType)
		//switch rType {
		//case 1:
		//	WSLogin(conn, rJson)
		//	if err != nil {
		//		fmt.Println("WSLogin err:", err.Error())
		//		continue
		//	}
		//case 2:
		//	WSDownload(conn, rJson)
		//	if err != nil {
		//		fmt.Println("WSDownload err:", err.Error())
		//		continue
		//	}
		//case 3:
		//	WSUpload(conn, rJson)
		//	if err != nil {
		//		fmt.Println("WSUpload err:", err.Error())
		//		continue
		//	}
		//}
	}
}
