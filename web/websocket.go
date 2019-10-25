package web

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"golang.org/x/net/websocket"
	"log"
	"redisgo/utils"
	"time"
)

var (
	connMap = make(map[string] *websocket.Conn)
	redisChanMap = make(map[string] map[string] chan int)
)

func WSHandler(conn *websocket.Conn) {
	seckey := conn.Request().Header.Get("Sec-Websocket-Key")
	log.Println("Websocket新建连接:", conn.RemoteAddr().String(), "=>", conn.LocalAddr().String(), seckey)
	connMap[seckey] = conn
	for {
		var reply string
		if err := websocket.Message.Receive(conn, &reply); err != nil {
			log.Println("Websocket连接断开:", err.Error())
			_ = conn.Close()
			return
		}
		rJson, err := simplejson.NewJson([]byte(reply))
		if err != nil {
			log.Println("接收到的值Json解析错误:", err.Error())
			return
		}
		rType, _ := rJson.Get("type").Int()
		switch rType {
		case 1: // 查询info信息
			ip, _ := rJson.Get("ip").String()
			container := utils.ContainerMap[ip]
			log.Println("收到查询info的命令, IP: " + ip)
			go func() {
				for {
					d, _ := json.Marshal(container.GetInfo())
					err := sendResponse(conn, 1, 0, ip, string(d))
					// 如果websocket断开, 退出协程
					if err != nil {
						return
					}
					time.Sleep(time.Second)
				}
			}()
		case 2:
			ip, _ := rJson.Get("ip").String()
			channel, _ := rJson.Get("channel").String()
			comm, _ := rJson.Get("command").String()
			container := utils.ContainerMap[ip]
			if comm == "open" {
				log.Println("收到订阅的命令, IP: " + ip + " Channel: " + channel)
				if redisChanMap[ip] == nil {
					redisChanMap[ip] = make(map[string] chan int)
				}
				if redisChanMap[ip][channel] == nil {
					redisChanMap[ip][channel] = make(chan int)
				}
				go func(command chan int) {
					pubsub := container.Subscribe(channel)

					// Wait for confirmation that subscription is created before publishing anything.
					_, err := pubsub.Receive()
					if err != nil {
						panic(err)
					}
					_ = sendResponse(conn, 2, 0, channel, "")

					// Go channel which receives messages.
					ch := pubsub.Channel()
					go func() {
						for msg := range ch {
							err := sendResponse(conn, 2, 1, msg.Channel, msg.Payload)
							if err != nil {
								_ = pubsub.Close()
								return
							}
						}
					}()
					for {
						select {
						case <- command:
							_ = pubsub.Close()
							_ = sendResponse(conn, 2, -1, ip, channel)
						default:
							// 发送心跳信号检测Websocket连接是否断开, 如果断开, 需要取消订阅该连接下的所有订阅
							err := sendResponse(conn, 0, 0, "", "")
							if err != nil {
								log.Println("Websocket连接断开, 取消订阅 IP: " + ip + " Channel: " + channel)
								_ = pubsub.Close()
								return
							}
							time.Sleep(time.Minute)
						}
					}
				}(redisChanMap[ip][channel])
			} else if comm == "close" {
				log.Println("收到取消订阅的命令, IP: " + ip + " Channel: " + channel)
				redisChanMap[ip][channel] <- 88
			}
		}
	}
}
