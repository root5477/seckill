package service

import (
	"github.com/gomodule/redigo/redis"
	"encoding/json"
	"secProxy/conf"
	"secProxy/log"
	"secProxy/model"
	"fmt"
)

func WriteHandle()  {
	//
	for {
		reqWithCookie := <- conf.SecReqChan
		coon := conf.RedisPoolProxy2Layer.Get()
		data, err := json.Marshal(reqWithCookie)
		if err != nil {
			log.Errorf("json marshal failed, err:%v", err)
			coon.Close()
			continue
		}
		_, err = coon.Do("LPUSH", "sec_queue", data)
		if err != nil {
			log.Errorf("lpush failed, err:%v, reqWithCookie:%v", err, *reqWithCookie)
			coon.Close()
			continue
		}
		coon.Close()
	}
}

func ReadHandle()  {
	for {
		conn := conf.RedisPoolProxy2Layer.Get()
		res, err := redis.String(conn.Do("BLPOP", conf.KillConf.Redis.Layer2ProxyQueueName, 0))
		if err != nil {
			log.Errorf("BLPOP from layer2proxy failed, err:%v, res:%v", err, res)
			conn.Close()
			continue
		}
		resp := &model.SecResponse{}
		err = json.Unmarshal([]byte(res), resp)
		if err != nil {
			log.Errorf("json unmarshal resp [%v] failed, err:%v", res, err)
			conn.Close()
			continue
		}
		userKey := fmt.Sprintf("%v-%v", resp.UserId, resp.ProductId)
		respChan, ok := conf.RespChanMap[userKey]
		if !ok {
			log.Errorf("resp chan not exist, respChan:%s, resp:%v", respChan, resp)
			conn.Close()
			continue
		}
		respChan <- resp
	}
}
