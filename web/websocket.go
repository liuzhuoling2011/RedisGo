package web

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"golang.org/x/net/websocket"
	"redisgo/utils"
	"time"
)

func WSHandler(conn *websocket.Conn) {
	fmt.Printf("Websocket新建连接: %s -> %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())

	for {
		var reply string
		if err := websocket.Message.Receive(conn, &reply); err != nil {
			fmt.Println("Websocket连接断开:", err.Error())
			conn.Close()
			for _, c := range utils.ContainerMap {
				c.Status = 0
			}
			return
		}
		rJson, err := simplejson.NewJson([]byte(reply))
		if err != nil {
			fmt.Println("receive err:", err.Error())
			return
		}
		rType, _ := rJson.Get("type").Int()
		switch rType {
		case 1:
			ip, _ := rJson.Get("ip").String()
			container := utils.ContainerMap[ip]
			if container.Status == 1 {
				continue
			}
			fmt.Println("收到查询info的命令, IP: " + ip)
			go func() {
				for {
					d, _ := json.Marshal(container.GetInfo())
					err := sendResponse(conn, 0, 0, ip, string(d))
					// 如果websocket断开, 退出协程
					if err != nil {
						return
					}
					time.Sleep(time.Second)
				}
			}()
			container.Status = 1
		}
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
