package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var KillConf *Config

type Server struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

type Db struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	PassWord string `json:"pass_word"`
	Database string `json:"database"`
}

type Redis struct {
	Addr        string `json:"addr"`
	Pwd         string `json:"pwd"`
	MaxIdle     int    `json:"max_idle"`
	MaxActive   int    `json:"max_active"`
	IdleTimeout int    `json:"idle_timeout"`
	DbOption    int    `json:"db_option"`
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
	Server    `json:"server"`
	Db        `json:"db"`
	Redis     `json:"redis"`
	Etcd      `json:"etcd"`
	Log       `json:"log"`
	SecretKey `json:"secret_key"`
	Limit     `json:"limit"`
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

func init() {
	KillConf = EnvInit(`/Users/chenqi/go/src/secProxy/conf/seckill.json`)
}
