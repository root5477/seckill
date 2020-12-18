package conf

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.etcd.io/etcd/clientv3"
	"secProxy/model"
	"sync"
	"time"
)

var RedisPool *redis.Pool
var RedisPoolProxy2Layer *redis.Pool

var CliEtcd *clientv3.Client

var SecInfosMap map[int] *model.SecInfo

var RWLockerOfSecInfo sync.RWMutex
var RWBlackIpLocker sync.RWMutex
var RWBlackIdLocker sync.RWMutex

//写入redis和读取redis 的chan
var SecReqChan chan *model.SecRequestWithCookie


func init()  {
	SecInfosMap = make(map[int] *model.SecInfo, 1024)
	SecReqChan = make(chan *model.SecRequestWithCookie, KillConf.Server.SecReqChanSize)
}

func InitRedis() (err error) {
	RedisPool = &redis.Pool{
		MaxIdle:KillConf.Redis.MaxIdle,
		MaxActive:KillConf.Redis.MaxActive,
		IdleTimeout:time.Duration(KillConf.Redis.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			dialPwd := redis.DialPassword(KillConf.Redis.Pwd)
			dbOption := redis.DialDatabase(KillConf.Redis.DbOption)
			Conns, err := redis.Dial("tcp", KillConf.Redis.Addr, dialPwd, dbOption)
			if err != nil {
				return nil, err
			}
			return Conns, nil
		},
	}
	c := RedisPool.Get()
	_, err = c.Do("ping")
	if err != nil {
		return err
	}
	return nil
}

func InitProxy2LayerRedis() (err error) {
	RedisPoolProxy2Layer = &redis.Pool{
		MaxIdle:KillConf.RedisProxy2layer.MaxIdle,
		MaxActive:KillConf.RedisProxy2layer.MaxActive,
		IdleTimeout:time.Duration(KillConf.RedisProxy2layer.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			dialPwd := redis.DialPassword(KillConf.RedisProxy2layer.Pwd)
			dbOption := redis.DialDatabase(KillConf.RedisProxy2layer.DbOption)
			Conns, err := redis.Dial("tcp", KillConf.RedisProxy2layer.Addr, dialPwd, dbOption)
			if err != nil {
				return nil, err
			}
			return Conns, nil
		},
	}
	c := RedisPoolProxy2Layer.Get()
	_, err = c.Do("ping")
	if err != nil {
		return err
	}
	return nil
}


func InitEtcd() (err error) {

	CliEtcd, err = clientv3.New(clientv3.Config{
		Endpoints:[]string {KillConf.Etcd.Addr},
		DialTimeout: time.Duration(KillConf.Etcd.TimeOut) * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err is:%v\n", err)
		return err
	}
	fmt.Println("connect to etcd success!")
	return
}






