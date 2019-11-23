package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type RedisInfo struct {
	// Server
	Redis_version     string `json:"redis_version"`
	Redis_git_sha1    string `json:"redis_git_sha_1"`
	Redis_git_dirty   string `json:"redis_git_dirty"`
	Redis_build_id    string `json:"redis_build_id"`
	Redis_mode        string `json:"redis_mode"`
	Os                string `json:"os"`
	Arch_bits         string `json:"arch_bits"`
	Multiplexing_api  string `json:"multiplexing_api"`
	Process_id        string `json:"process_id"`
	Run_id            string `json:"run_id"`
	Tcp_port          string `json:"tcp_port"`
	Uptime_in_seconds string `json:"uptime_in_seconds"`
	Uptime_in_days    string `json:"uptime_in_days"`
	Hz                string `json:"hz"`
	Lru_clock         string `json:"lru_clock"`
	Executable        string `json:"executable"`
	Config_file        string `json:"config_file"`
	//Client
	Connected_clients          string `json:"connected_clients"`
	Client_longest_output_list string `json:"client_longest_output_list"`
	Client_biggest_input_buf   string `json:"client_biggest_input_buf"`
	Blocked_clients            string `json:"blocked_clients"`
	//Memory
	Used_memory               string `json:"used_memory"`
	Used_memory_human         string `json:"used_memory_human"`
	Used_memory_rss           string `json:"used_memory_rss"`
	Used_memory_rss_human     string `json:"used_memory_rss_human"`
	Used_memory_peak          string `json:"used_memory_peak"`
	Used_memory_peak_human    string `json:"used_memory_peak_human"`
	Total_system_memory       string `json:"total_system_memory"`
	Total_system_memory_human string `json:"total_system_memory_human"`
	Used_memory_lua           string `json:"used_memory_lua"`
	Used_memory_lua_human     string `json:"used_memory_lua_human"`
	Maxmemory                 string `json:"maxmemory"`
	Maxmemory_human           string `json:"maxmemory_human"`
	Maxmemory_policy          string `json:"maxmemory_policy"`
	Mem_fragmentation_ratio   string `json:"mem_fragmentation_ratio"`
	Mem_allocator             string `json:"mem_allocator"`
	// Persistence
	Loading                      string `json:"loading"`
	Rdb_changes_since_last_save  string `json:"rdb_changes_since_last_save"`
	Rdb_bgsave_in_progress       string `json:"rdb_bgsave_in_progress"`
	Rdb_last_save_time           string `json:"rdb_last_save_time"`
	Rdb_last_bgsave_status       string `json:"rdb_last_bgsave_status"`
	Rdb_last_bgsave_time_sec     string `json:"rdb_last_bgsave_time_sec"`
	Rdb_current_bgsave_time_sec  string `json:"rdb_current_bgsave_time_sec"`
	Aof_enabled                  string `json:"aof_enabled"`
	Aof_current_size             string `json:"aof_current_size"`
	Aof_rewrite_in_progress      string `json:"aof_rewrite_in_progress"`
	Aof_rewrite_secheduled       string `json:"aof_rewrite_secheduled"`
	Aof_last_rewrite_time_sec    string `json:"aof_last_rewrite_time_sec"`
	Aof_current_rewrite_time_sec string `json:"aof_current_rewrite_time_sec"`
	Aof_last_bgrewrite_status    string `json:"aof_last_bgrewrite_status"`
	Aof_last_write_status        string `json:"aof_last_write_status"`
	//Stats
	Total_connections_received string `json:"total_connections_received"`
	Total_commands_processed   string `json:"total_commands_processed"`
	Instantaneous_ops_per_sec  string `json:"instantaneous_ops_per_sec"`
	Total_net_input_bytes      string `json:"total_net_input_bytes"`
	Total_net_output_bytes     string `json:"total_net_output_bytes"`
	Instantaneous_input_kbps   string `json:"instantaneous_input_kbps"`
	Instantaneous_output_kbps  string `json:"instantaneous_output_kbps"`
	Rejected_connections       string `json:"rejected_connections"`
	Sync_full                  string `json:"sync_full"`
	Sync_partial_ok            string `json:"sync_partial_ok"`
	Sync_partial_err           string `json:"sync_partial_err"`
	Expired_keys               string `json:"expired_keys"`
	Evicted_keys               string `json:"evicted_keys"`
	Keyspace_hits              string `json:"keyspace_hits"`
	Keyspace_misses            string `json:"keyspace_misses"`
	Pubsub_channels            string `json:"pubsub_channels"`
	Pubsub_patterns            string `json:"pubsub_patterns"`
	Latest_fork_usec           string `json:"latest_fork_usec"`
	Migrate_cached_sockets     string `json:"migrate_cached_sockets"`

	//Replication
	Role                           string `json:"role"`
	Connected_slaves               string `json:"connected_slaves"`
	Master_repl_offset             string `json:"master_repl_offset"`
	Repl_backlog_active            string `json:"repl_backlog_active"`
	Repl_backlog_size              string `json:"repl_backlog_size"`
	Repl_backlog_first_byte_offset string `json:"repl_backlog_first_byte_offset"`
	Repl_backlog_histlen           string `json:"repl_backlog_histlen"`

	// Cpu
	Used_cpu_sys           string `json:"used_cpu_sys"`
	Used_cpu_user          string `json:"used_cpu_user"`
	Used_cpu_sys_children  string `json:"used_cpu_sys_children"`
	Used_cpu_user_children string `json:"used_cpu_user_children"`

	// Cluster
	Cluster_enabled string `json:"cluster_enabled"`

	// Keyspace
	Db0 string `json:"db0"`
	Db1 string `json:"db1"`
	Db2 string `json:"db2"`
	Db3 string `json:"db3"`
	Db4 string `json:"db4"`
	Db5 string `json:"db5"`
	Db6 string `json:"db6"`
	Db7 string `json:"db7"`
	Db8 string `json:"db8"`
	Db9 string `json:"db9"`
	Db10 string `json:"db10"`
	Db11 string `json:"db11"`
	Db12 string `json:"db12"`
	Db13 string `json:"db13"`
	Db14 string `json:"db14"`
	Db15 string `json:"db15"`
}

type RedisLog struct {
	Id          int64 `json:"id"`
	Time        int64 `json:"time"`
	Time_used   int64 `json:"time_used"`
	Msg         string `json:"msg"`
}

type RedisLogList []*RedisLog

type RedisClient struct {
	Id          string `json:"id"`
	Addr        string `json:"addr"`
	Fd          string `json:"fd"`
	Name        string `json:"name"`
	Age         string `json:"age"`
	Idle        string `json:"idle"`
	Flag        string `json:"flag"`
	Db          string `json:"db"`
	Sub         string `json:"sub"`
	Psub        string `json:"psub"`
	Multi       string `json:"multi"`
	Qbuf        string `json:"qbuf"`
	Qbuf_free   string `json:"qbuf-free"`
	Obl         string `json:"obl"`
	Oll         string `json:"oll"`
	Omem        string `json:"omem"`
	Events      string `json:"events"`
	Cmd         string `json:"cmd"`
}

type RedisClientList []*RedisClient

type Config struct {
	Ip       string `json:"ip"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Db       int    `json:"db"`
	Name     string `json:"name"`
}

type Container struct {
	Config
	Status     uint8  `json:"status"`
	redis      *redis.Client
}

var (
	ContainerMap = make(map[string]*Container)
	RedisConfigs []Config
)

func InitConfig() bool {
	filePtr, err := os.Open("./config.json")
	defer filePtr.Close()
	if err != nil { return true }

	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&RedisConfigs)
	if err != nil {
		log.Println("配置文件解析失败", err.Error())
		return false
	}
	return true
}

func SaveConfig() bool {
	filePtr, err := os.OpenFile("./config.json", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	defer filePtr.Close()
	data, _ := json.MarshalIndent(RedisConfigs, "", "    ")
	_, err = filePtr.Write(data)
	if err != nil {
		fmt.Println("保存配置文件失败", err.Error())
		return false
	}
	return true
}

func InitContainers() bool {
	for _, config := range RedisConfigs {
		AddContainer(config)
	}
	return true
}

func AddContainer(config Config) bool {
	if ContainerMap[config.Ip] != nil {
		log.Println("redis ip 重复, 请检查", config.Ip)
		return false
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.Ip + ":" + strconv.Itoa(config.Port),
		Password: config.Password,
		DB:       config.Db,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Println("redis连接错误", config.Ip, err.Error())
		return false
	}
	container := &Container{config, 0, client}
	ContainerMap[config.Ip] = container
	return true
}

func UpdateContainer(config Config) bool {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Ip + ":" + strconv.Itoa(config.Port),
		Password: config.Password,
		DB:       config.Db,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Println("redis连接错误", config.Ip, err.Error())
		return false
	}
	ContainerMap[config.Ip].redis = client
	return true
}

func DeleteContainer(ip string) {
	delete(ContainerMap, ip)
	index := 0
	for ; index < len(RedisConfigs); {
		if RedisConfigs[index].Ip == ip {
			RedisConfigs = append(RedisConfigs[:index], RedisConfigs[index+1:]...)
			continue
		}
		index++
	}
	SaveConfig()
}

func (c *Container) GetInfo() *RedisInfo {
	infos := strings.Split(c.redis.Info().String(), "\r\n")
	rinfo := &RedisInfo{}
	infom := make(map[string]interface{})
	for _, line := range infos {
		args := strings.Split(line, ":")
		if len(args) < 2 {
			continue
		}
		infom[args[0]] = args[1]
	}
	b, _ := json.Marshal(infom)
	_ = json.Unmarshal(b, &rinfo)
	return rinfo
}

func (c *Container) GetLog() *RedisLogList {
	rll := RedisLogList{}
	logs, _ := c.redis.Do("SLOWLOG", "GET", "120").Result()
	if logs_interfaces, ok := logs.([]interface{}); ok {
		for _, logsi := range logs_interfaces {
			if log, ok := logsi.([]interface{}); ok {
				rlog := &RedisLog{}
				rlog.Id = log[0].(int64)
				rlog.Time = log[1].(int64)
				rlog.Time_used = log[2].(int64)
				rlog.Msg = fmt.Sprintf("%s", log[3])
				rll = append(rll, rlog)
			}
		}
	}
	return &rll
}

func (c *Container) GetClients() *RedisClientList {
	rcl := RedisClientList{}
	infos := strings.Split(c.redis.ClientList().String(), "\n")
	for _, line := range infos {
		if line == "" {
			continue
		}
		rclient := &RedisClient{}
		clientm := make(map[string]interface{})
		args := strings.Split(line, " ")
		for _, pair := range args {
			if pair == "client" || pair == "list:" {
				continue
			}
			pairs := strings.Split(pair, "=")
			clientm[pairs[0]] = pairs[1]
		}
		b, _ := json.Marshal(clientm)
		_ =json.Unmarshal(b, &rclient)
		rcl = append(rcl, rclient)
	}
	return &rcl
}

func (c *Container) PublishMsg(key string, msg string) {
	info := c.redis.Publish(key, msg)
	fmt.Println(info)
}

func (c *Container) Subscribe(channel string) *redis.PubSub {
	channels := strings.Split(channel, " ")
	return c.redis.PSubscribe(channels...)
}

func (c *Container) Execute(command string) interface{} {
	args := strings.Split(command, " ")
	var commands = make([]interface{}, len(args))
	for i, v := range args {
		commands[i] = v
	}
	info, _ := c.redis.Do(commands...).Result()
	if info == nil {
		return "执行错误, 请检查输入的Redis命令"
	}
	return info
}

func(c *Container) Rename(key string, newName string) string {
	t, _ := c.redis.Rename(key, newName).Result()
	return t
}

func(c *Container) GetType(key string) string {
	t, _ := c.redis.Type(key).Result()
	return t
}

func(c *Container) GetTTL(key string) time.Duration {
	t, _ := c.redis.TTL(key).Result()
	return t
}

func(c *Container) SetTTL(key string, ttl int) bool {
	var ret bool
	if ttl == -1 {
		ret, _ = c.redis.Persist(key).Result()
	} else {
		var seconds int = int(time.Second)
		ret, _ = c.redis.Expire(key, time.Duration(ttl * seconds)).Result()
	}
	return ret
}

func(c *Container) DeleteKey(key string) int64 {
	t, _ := c.redis.Del(key).Result()
	return t
}

func(c *Container) GetLen(key string) int64 {
	var ret int64
	switch c.GetType(key) {
	case "string":
		ret, _ = c.redis.StrLen(key).Result()
	case "list":
		ret, _ = c.redis.LLen(key).Result()
	case "hash":
		ret, _ = c.redis.HLen(key).Result()
	case "set":
		ret, _ = c.redis.SCard(key).Result()
	case "zset":
		ret, _ = c.redis.ZCard(key).Result()
	}
	return ret
}

func (c *Container) SelectDB(db string) interface{} {
	info, _ := c.redis.Do("select", db).Result()
	return info
}

type KeyStruct struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Len        int64  `json:"len"`
	TTL        time.Duration `json:"ttl"`
}

type KeyScanStruct struct {
	Cursor     uint64 `json:"cursor"`
	Keys       []KeyStruct `json:"keys"`
}

func (c *Container) ScanKeys(cursor int, match string, count int) interface{} {
	var keylist []KeyStruct
	scursor := cursor
	var scount = 0
	for {
		info, cur, _ := c.redis.Scan(uint64(scursor), match, int64(count)).Result()
		for _, v := range info {
			keylist = append(keylist, KeyStruct{Name: v, Type: c.GetType(v), TTL: c.GetTTL(v), Len: c.GetLen(v)})
		}
		scount += len(info)
		scursor = int(cur)
		if scount >= count || scursor == 0 {
			break
		}
	}

	return KeyScanStruct{Cursor: uint64(scursor), Keys: keylist}
}

func (c *Container) GetKeys(command string) interface{} {
	var keylist []KeyStruct
	info, _ := c.redis.Keys(command).Result()
	for _, v := range info {
		keylist = append(keylist, KeyStruct{Name: v, Type: c.GetType(v), TTL: c.GetTTL(v), Len: c.GetLen(v)})
	}
	return keylist
}

func (c *Container) GetStringValue(key string) string {
	info, _ := c.redis.Get(key).Result()
	return info
}

func (c *Container) SetStringValue(key string, value string, ttl int) string {
	info, _ := c.redis.Set(key, value, time.Duration(ttl)).Result()
	return info
}

func (c *Container) GetListValueAll(key string) []string {
	info, _ := c.redis.LRange(key, 0, -1).Result()
	return info
}

func (c *Container) GetListValueIndex(key string, pos int) string {
	info, _ := c.redis.LIndex(key, int64(pos)).Result()
	return info
}

func (c *Container) SetListValue(key string, value string, pos int) string {
	info, _ := c.redis.LSet(key, int64(pos), value).Result()
	return info
}

func (c *Container) DeleteListValue(key string, pos int) int64 {
	v := c.GetListValueIndex(key, pos)
	value := md5V(string(pos) + v)
	c.SetListValue(key, value, pos)
	info, _ := c.redis.LRem(key, 1, value).Result()
	return info
}

func (c *Container) PushListValue(key string, value string, pos int) int64 {
	var info int64
	if pos == 0 {
		info, _ = c.redis.LPush(key, value).Result()
	} else if pos == -1{
		info, _ = c.redis.RPush(key, value).Result()
	}
	return info
}

func (c *Container) GetHashValueAll(key string) map[string]string {
	info, _ := c.redis.HGetAll(key).Result()
	return info
}

type HashScanStruct struct {
	Cursor     uint64 `json:"cursor"`
	Keys       []string `json:"keys"`
}

func (c *Container) ScanHashValue(key string, cursor int, match string, count int) HashScanStruct {
	var keylist []string
	scursor := cursor
	var scount = 0
	for {
		info, cur, _ := c.redis.HScan(key, uint64(scursor), match, int64(count)).Result()
		for _, v := range info {
			keylist = append(keylist, v)
		}
		scount += len(info) / 2
		scursor = int(cur)
		if scount >= count || scursor == 0 {
			break
		}
	}

	return HashScanStruct{Cursor: uint64(scursor), Keys: keylist}
}

func (c *Container) GetSetValueAll(key string) []string {
	info, _ := c.redis.SMembers(key).Result()
	return info
}

//func (c *Container) ScanSetValue(key string, cursor int, match string, count int) []string {
//	info, cur, _ := c.redis.SScan(key, uint64(cursor), match, int64(count)).Result()
//	return info
//}

//func (c *Container) GetZSetValueAll(key string) []string {
//	info, _ := c.redis.SMembers(key).Result()
//	return info
//}