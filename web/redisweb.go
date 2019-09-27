package web

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

func StartServer(port uint, access bool) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("invalid port %d", port)
	}

	http.HandleFunc("/", RootHandle)

	http.Handle("/ws", websocket.Handler(WSHandler))
	if access {
		return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}
	fmt.Println("现在只监听localhost，请注意")
	return http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}
