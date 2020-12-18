package service

import (
	"encoding/json"
	"secProxy/conf"
	"secProxy/log"
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

}
