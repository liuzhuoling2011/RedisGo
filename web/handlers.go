package web

import (
	"net/http"
	"redisgo/utils"
	"strconv"
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
	case "logs":
		ip := r.Form.Get("ip")
		container := utils.ContainerMap[ip]
		sendHttpResponse(w, "", container.GetLog())
	case "clients":
		ip := r.Form.Get("ip")
		container := utils.ContainerMap[ip]
		sendHttpResponse(w, "", container.GetClients())
	case "delete":
		ip := r.Form.Get("ip")
		utils.DeleteContainer(ip)
		sendHttpResponse(w, "删除成功", "")
	case "edit":
		ip := r.Form.Get("ip")
		name := r.Form.Get("name")
		password := r.Form.Get("password")
		port, _ := strconv.Atoi(r.Form.Get("port"))
		db, _ := strconv.Atoi(r.Form.Get("db"))
		container := utils.ContainerMap[ip]
		if container.Password != password || container.Port != port || container.Db != db {
			if !utils.UpdateContainer(utils.Config{Ip:ip, Name:name, Password:password, Port:port, Db:db}) {
				sendHttpErrorResponse(w, -1, "修改错误, 请检查redis配置")
				return
			}
		} else {
			container.Name = name
		}
		sendHttpResponse(w, "修改成功", utils.ContainerMap[ip])
		utils.SaveConfig()
	case "add":
		ip := r.Form.Get("ip")
		name := r.Form.Get("name")
		password := r.Form.Get("password")
		port, _ := strconv.Atoi(r.Form.Get("port"))
		db, _ := strconv.Atoi(r.Form.Get("db"))
		if utils.AddContainer(utils.Config{Ip:ip, Name:name, Password:password, Port:port, Db:db}) {
			sendHttpResponse(w, "添加成功", utils.ContainerMap[ip])
			utils.SaveConfig()
		} else {
			sendHttpErrorResponse(w, -1, "添加错误, 请检查redis配置是否重复或者正确")
		}
	}
}
