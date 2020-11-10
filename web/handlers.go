package web

import (
	"io/ioutil"
	"net/http"
	"redisgo/utils"
	"runtime"
	"strconv"
)

const (
	cmdList   = "list"
	cmdInfo   = "info"
	cmdLog    = "logs"
	cmdDelete = "delete"
	cmdSet    = "set"
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
	case cmdList:
		sendHttpResponse(w, "", utils.ContainerMap)
	case cmdInfo:
		ip := r.Form.Get("ip")
		container := utils.ContainerMap[ip]
		sendHttpResponse(w, "", container.GetInfo())
	case cmdLog:
		ip := r.Form.Get("ip")
		container := utils.ContainerMap[ip]
		sendHttpResponse(w, "", container.GetLog())
	case "clients":
		ip := r.Form.Get("ip")
		container := utils.ContainerMap[ip]
		sendHttpResponse(w, "", container.GetClients())
	case cmdDelete:
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
			if !utils.UpdateContainer(utils.Config{Ip: ip, Name: name, Password: password, Port: port, Db: db}) {
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
		if utils.AddContainer(utils.Config{Ip: ip, Name: name, Password: password, Port: port, Db: db}) {
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

func DataHandle(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	ip := r.Form.Get("ip")
	key := r.Form.Get("key")
	container := utils.ContainerMap[ip]
	method := r.Form.Get("method")
	switch method {
	case "select_db":
		db := r.Form.Get("db")
		sendHttpResponse(w, "", container.SelectDB(db))
	case "get_keys":
		cursor, _ := strconv.Atoi(r.Form.Get("cursor"))
		match := r.Form.Get("match")
		count, _ := strconv.Atoi(r.Form.Get("count"))
		sendHttpResponse(w, "", container.ScanKeys(cursor, match, count))
	case "get_key_value":
		rtype := r.Form.Get("type")
		if rtype == "string" {
			sendHttpResponse(w, "", container.GetStringValue(key))
		} else if rtype == "list" {
			start, _ := strconv.Atoi(r.Form.Get("start"))
			end, _ := strconv.Atoi(r.Form.Get("end"))
			sendHttpResponse(w, "", container.GetListValueRange(key, start, end))
		} else if rtype == "hash" {
			cursor, _ := strconv.Atoi(r.Form.Get("cursor"))
			match := r.Form.Get("match")
			count, _ := strconv.Atoi(r.Form.Get("count"))
			sendHttpResponse(w, "", container.ScanHashValue(key, cursor, match, count))
		} else if rtype == cmdSet {
			cursor, _ := strconv.Atoi(r.Form.Get("cursor"))
			match := r.Form.Get("match")
			count, _ := strconv.Atoi(r.Form.Get("count"))
			sendHttpResponse(w, "", container.ScanSetValue(key, cursor, match, count))
		} else if rtype == "zset" {
			cursor, _ := strconv.Atoi(r.Form.Get("cursor"))
			match := r.Form.Get("match")
			count, _ := strconv.Atoi(r.Form.Get("count"))
			sendHttpResponse(w, "", container.ScanZSetValue(key, cursor, match, count))
		}
	case "rm_key":
		sendHttpResponse(w, "", container.DeleteKey(key))
	case "update_ttl":
		ttl, _ := strconv.Atoi(r.Form.Get("ttl"))
		sendHttpResponse(w, "", container.SetTTL(key, ttl))
	case "get_ttl":
		sendHttpResponse(w, "", container.GetTTL(key))
	case "rename":
		newName := r.Form.Get("new_name")
		sendHttpResponse(w, "", container.Rename(key, newName))
	case "string_ops":
		ops := r.Form.Get("ops")
		if ops == cmdSet {
			ttl, _ := strconv.Atoi(r.Form.Get("ttl"))
			value := r.Form.Get("value")
			sendHttpResponse(w, "", container.SetStringValue(key, value, ttl))
		}
	case "list_ops":
		ops := r.Form.Get("ops")
		if ops == "push" {
			pos, _ := strconv.Atoi(r.Form.Get("pos"))
			value := r.Form.Get("value")
			sendHttpResponse(w, "", container.PushListValue(key, value, pos))
		} else if ops == cmdDelete {
			pos, _ := strconv.Atoi(r.Form.Get("pos"))
			sendHttpResponse(w, "", container.DeleteListValue(key, pos))
		} else if ops == cmdSet {
			pos, _ := strconv.Atoi(r.Form.Get("pos"))
			value := r.Form.Get("value")
			sendHttpResponse(w, "", container.SetListValue(key, pos, value))
		}
	case "hash_ops":
		ops := r.Form.Get("ops")
		hashKey := r.Form.Get("hash_key")
		if ops == cmdDelete {
			sendHttpResponse(w, "", container.DeleteHashValue(key, hashKey))
		} else if ops == cmdSet {
			value := r.Form.Get("value")
			sendHttpResponse(w, "", container.SetHashValue(key, hashKey, value))
		}
	case "set_ops":
		ops := r.Form.Get("ops")
		setKey := r.Form.Get("set_key")
		if ops == cmdDelete {
			sendHttpResponse(w, "", container.DeleteSetValue(key, setKey))
		} else if ops == cmdSet {
			value := r.Form.Get("value")
			sendHttpResponse(w, "", container.SetSetValue(key, setKey, value))
		}
	case "zset_ops":
		ops := r.Form.Get("ops")
		zsetKey := r.Form.Get("zset_key")
		if ops == cmdDelete {
			sendHttpResponse(w, "", container.DeleteZSetValue(key, zsetKey))
		} else if ops == cmdSet {
			value, _ := strconv.ParseFloat(r.Form.Get("value"), 64)
			sendHttpResponse(w, "", container.SetZSetValue(key, zsetKey, value))
		}
	}
}
