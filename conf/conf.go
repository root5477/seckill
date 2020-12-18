package conf

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"time"
)

var KillConf *Config

type Server struct {
	Addr                         string `json:"addr"`
	Port                         int    `json:"port"`
	WriteProxy2LayerGoroutineNum int    `json:"write_proxy_2_layer_goroutine_num"`
	ReadLayer2ProxyGoroutineNum  int    `json:"read_layer_2_proxy_goroutine_num"`
	SecReqChanSize               int    `json:"sec_req_chan_size"`
}

type Db struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	PassWord string `json:"pass_word"`
	Database string `json:"database"`
}

//配置redis，如：ip黑名单/id黑名单
type Redis struct {
	Addr                 string `json:"addr"`
	Pwd                  string `json:"pwd"`
	MaxIdle              int    `json:"max_idle"`
	MaxActive            int    `json:"max_active"`
	IdleTimeout          int    `json:"idle_timeout"`
	DbOption             int    `json:"db_option"`
	Layer2ProxyQueueName string `json:"layer_2_proxy_queue_name"`
}

//接入层redis -> 业务逻辑
type RedisProxy2layer struct {
	Addr                 string `json:"addr"`
	Pwd                  string `json:"pwd"`
	MaxIdle              int    `json:"max_idle"`
	MaxActive            int    `json:"max_active"`
	IdleTimeout          int    `json:"idle_timeout"`
	DbOption             int    `json:"db_option"`
	Proxy2LayerQueueName string `json:"proxy_2_layer_queue_name"`
}

type Etcd struct {
	Addr    string `json:"addr"`
	Pwd     string `json:"pwd"`
	TimeOut int    `json:"time_out"`
	SecKey  string `json:"sec_key"`
}

type Log struct {
	MinLevel string `json:"min_level"`
}

type SecretKey struct {
	CookieSecretKey string `json:"cookie_secret_key"`
}

type Limit struct {
	UserSecAccessLimit int      `json:"user_sec_access_limit"`
	IpSecAccessLimit   int      `json:"ip_sec_access_limit"`
	ReferWhiteList     []string `json:"refer_white_list"`
}

type Config struct {
	Server                     `json:"server"`
	Db                         `json:"db"`
	Redis                      `json:"redis"`
	RedisProxy2layer           `json:"redis_proxy2layer"`
	Etcd                       `json:"etcd"`
	Log                        `json:"log"`
	SecretKey                  `json:"secret_key"`
	Limit                      `json:"limit"`
	IpBlackMap map[string]bool `json:"ip_black_map"`
	IdBlackMap map[string]bool `json:"id_black_map"`
}

func EnvInit(confPath string) *Config {
	//InitLogger()
	Config := &Config{}
	confByte, err := ioutil.ReadFile(confPath)
	if err != nil {
		fmt.Printf("读取配置文件%s失败, err is:%v", confPath, err)
		panic(err)
	}
	errUnmarshal := json.Unmarshal(confByte, Config)
	if errUnmarshal != nil {
		fmt.Printf("json.Unmarshal失败, err is:%v", errUnmarshal)
		panic(errUnmarshal)
	}
	return Config
}

//加载黑名单数据：ip & id
func LoadBlackList() (err error) {
	c := RedisPool.Get()
	defer c.Close()
	reply, err := redis.Strings(c.Do("hgetall", "id_black_list"))
	if err != nil {
		return
	}
	for _, v := range reply {
		KillConf.IdBlackMap[v] = true
	}

	reply, err = redis.Strings(c.Do("hgetall", "ip_black_list"))
	if err != nil {
		return
	}
	for _, v := range reply {
		KillConf.IpBlackMap[v] = true
	}

	go SyncIdBlackList()
	go SyncIpBlackList()
	return
}

func SyncIpBlackList() {
	for {
		conn := RedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "black_ip_list", time.Second)
		ip, err := redis.String(reply, err)
		if err != nil {
			continue
		}
		RWBlackIpLocker.Lock()
		KillConf.IpBlackMap[ip] = true
		RWBlackIpLocker.Unlock()
	}
}

func SyncIdBlackList() {
	for {
		conn := RedisPool.Get()
		defer conn.Close()
		reply, err := conn.Do("BLPOP", "black_id_list", time.Second)
		id, err := redis.String(reply, err)
		if err != nil {
			continue
		}
		RWBlackIdLocker.Lock()
		KillConf.IdBlackMap[id] = true
		RWBlackIdLocker.Unlock()
	}
}

func init() {
	KillConf = EnvInit(`./conf/seckill.json`)

}
