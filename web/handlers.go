package web

import (
	"net/http"
	"redisgo/utils"
)

func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		next.ServeHTTP(w, r)
	}
}

func RootHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		indexPage(w, r)
	} else {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(404)
	}
}

func ContainerHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	method := r.Form.Get("method")
	switch method {
	case "list":
		sendHttpResponse(w, "", utils.ContainerMap)
	case "info":
		ip := r.Form.Get("ip")
		container := utils.ContainerMap[ip]
		sendHttpResponse(w, "", container.GetInfo())
	}
}
