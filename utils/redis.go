package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"os"
)

type RicliInfo struct {
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
	Aof_enable                   string `json:"aof_enable"`
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
}

var containerMap = make(map[string]*Container)

type Config struct {
	Ip       string `json:"ip"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Db       int    `json:"db"`
	Name     string `json:"name"`
}

type Container struct {
	RemoteAddr string
	Name       string
	redis      *redis.Client
	Status     uint8
}

func init_config() {
	filePtr, err := os.Open("./config.json")
	defer filePtr.Close()
	if err != nil {
		fmt.Println("不存在默认配置文件, 将创建默认配置文件")
		filePtr1, _ := os.Create("./config.json")
		defer filePtr1.Close()
		defaultConfig := []Config{{"localhost", "", 6379, 0, "default"}}
		data, _ := json.MarshalIndent(defaultConfig, "", "    ")
		filePtr1.Write(data)
		return false
	}

	var configs []Config
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&configs)
	if err != nil {
		fmt.Println("配置文件解析失败", err.Error())
		return false
	}
	fmt.Println(configs)

	return true
}

//func (c *Container) interval() int64 {
//	now := time.Now().Unix()
//	i := now - c.LastTime.Unix()
//	return i
//}
//
//func (c *Container) getRicliInfo(remoteAddr string) *RicliInfo {
//	return c.RicliInfo
//}

//func getRicliInfo(remoteAddr string) *RicliInfo {
//	container := containerMap[remoteAddr]
//	if container != nil {
//		return container.RicliInfo
//	}
//	return nil
//}

//func addContainer(remoteAddr string, ricliInfo *RicliInfo) {
//	if r := containerMap[remoteAddr]; r != nil {
//		r.LastTime = time.Now()
//		r.RicliInfo = ricliInfo
//		r.Status = NORMAL
//	} else {
//		container := &Container{remoteAddr, remoteAddr,ricliInfo, time.Now(), NORMAL}
//		containerMap[remoteAddr] = container
//	}
//}

func GetContainers() map[string]*Container {
	return containerMap
}

func GetContainer(remoteAddr string) *Container {
	return containerMap[remoteAddr]
}
