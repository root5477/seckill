package main

import (
	"github.com/gin-gonic/gin"
	"secProxy/conf"
	"secProxy/dao"
	"secProxy/handle"
	"secProxy/log"
)

//遇到因为grpc或clientv3 biuld失败，切换为root，go mod init, go build main.go,
//然后修改go.mod replace google.golang.org/grpc => google.golang.org/grpc v1.26.0，再gobuild
func main()  {
	//init log
	err := log.InitLog("./log/seelog.xml", conf.KillConf)
	if err != nil {
		panic(err)
	}
	log.Debugf("InitLog success! ")
	log.Debugf("KillConf:[%v]", *conf.KillConf)
	//init redis
	err = conf.InitRedis()
	if err != nil {
		log.Errorf("InitRedis failed, err is:%v", err)
		panic(err)
	}
	log.Debugf("InitRedis success!")

	//load sec info
	err = dao.LoadSecInfoConf()
	if err != nil {
		log.Errorf("LoadSecInfoConf from etcd failed, err is:%v", err)
		panic(err)
	}
	log.Debugf("LoadSecInfoConf from redis success!")
	//init etcd watcher
	dao.InitSecProductWatcher()

	router := gin.Default()
	router.GET("/seckill", handle.SecKillHandle)
	router.GET("/secinfo", handle.SecInfo)

	router.Run(":8900")
}
