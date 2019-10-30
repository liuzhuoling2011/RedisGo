package web

import (
	"io/ioutil"
	"net/http"
	"redisgo/utils"
	"runtime"
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
	_ = r.ParseForm()
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
		if container == nil {
			sendHttpErrorResponse(w, -1, "修改错误, IP不存在")
			return
		}
		if container.Password != password || container.Port != port || container.Db != db {
			if !utils.UpdateContainer(utils.Config{Ip:ip, Name:name, Password:password, Port:port, Db:db}) {
				sendHttpErrorResponse(w, -1, "修改错误, 请检查redis配置")
				return
			}
		} else {
			container.Name = name
		}
		sendHttpResponse(w, "修改成功", utils.ContainerMap[ip])
		for index, conf := range utils.RedisConfigs {
			if conf.Ip == ip {
				utils.RedisConfigs[index].Name = name
				utils.RedisConfigs[index].Password = password
				utils.RedisConfigs[index].Port = port
				utils.RedisConfigs[index].Db = db
			}
		}
		utils.SaveConfig()
	case "add":
		ip := r.Form.Get("ip")
		name := r.Form.Get("name")
		password := r.Form.Get("password")
		port, _ := strconv.Atoi(r.Form.Get("port"))
		db, _ := strconv.Atoi(r.Form.Get("db"))
		if utils.AddContainer(utils.Config{Ip:ip, Name:name, Password:password, Port:port, Db:db}) {
			sendHttpResponse(w, "添加成功", utils.ContainerMap[ip])
			utils.RedisConfigs = append(utils.RedisConfigs, utils.Config{ip, password, port, db, name})
			utils.SaveConfig()
		} else {
			sendHttpErrorResponse(w, -1, "添加错误, 请检查redis配置是否重复或者正确")
		}
	case "publish":
		ip := r.Form.Get("ip")
		key := r.Form.Get("key")
		msg := r.Form.Get("msg")
		container := utils.ContainerMap[ip]
		container.PublishMsg(key, msg)
		sendHttpResponse(w, key, msg)
	case "execute":
		ip := r.Form.Get("ip")
		command := r.Form.Get("command")
		container := utils.ContainerMap[ip]
		sendHttpResponse(w, "", container.Execute(command))
	}
}

func SystemHandle(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	method := r.Form.Get("method")
	switch method {
	case "update":
		url := "http://www.zoranjojo.top:9925/api/v1/update?goos=" + runtime.GOOS + "&goarch=" + runtime.GOARCH
		url += "&version=" + utils.GetVersion() + "&prefix=" + utils.GetName()
		resp, err := http.Get(url)
		if err != nil {
			sendHttpErrorResponse(w, -1, "查找版本更新失败")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			sendHttpErrorResponse(w, -2, "查找版本更新失败")
		}
		sendHttpResponse(w, "", string(body))
	case "notice":
		url := "http://www.zoranjojo.top:9925/api/v1/notice?product=" + utils.GetName()
		resp, err := http.Get(url)
		if err != nil {
			sendHttpErrorResponse(w, -1, "查找通知信息失败")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			sendHttpErrorResponse(w, -2, "查找通知信息失败")
		}

		sendHttpResponse(w, "", string(body))
	}
}