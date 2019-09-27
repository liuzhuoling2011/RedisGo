package web

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"redisgo/utils"
)

var (
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
)

func init() {
	fmt.Println("handler init")

}

func RootHandle(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	_, err := client.Ping().Result()
	if err != nil {
		sendHttpErrorResponse(w, -1, err.Error())
		return
	}
	res := client.Info()
	fmt.Println(res)
	sendHttpResponse(w, "", res)
}
