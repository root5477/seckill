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
var CliEtcd *clientv3.Client
var SecInfosMap map[int] *model.SecInfo
var RWLockerOfSecInfo sync.RWMutex

func init()  {
	SecInfosMap = make(map[int] *model.SecInfo, 1024)
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






