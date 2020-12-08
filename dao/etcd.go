package dao

import (
	"go.etcd.io/etcd/clientv3"
	"context"
	"encoding/json"
	"fmt"
	."secProxy/conf"
	"secProxy/log"
	"secProxy/model"

	"time"

)

//加载秒杀商品信息
func LoadSecInfoConf() (err error) {
	err = InitEtcd()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	resp, err := CliEtcd.Get(ctx, "sec_info")
	cancel()
	if err != nil {
		return err
	}
	var SecInfos []model.SecInfo
	for _, v := range resp.Kvs {
		err = json.Unmarshal([]byte(v.Value), &SecInfos)
		if err != nil {
			return err
		}
	}
	fmt.Println("secInfos:", SecInfos)
	RWLockerOfSecInfo.Lock()
	for i := range SecInfos {
		SecInfosMap[SecInfos[i].ProductId] = &SecInfos[i]
	}
	RWLockerOfSecInfo.Unlock()
	return nil
}

func InitSecProductWatcher()  {
	go WatchSecProductInfoKey(KillConf.Etcd.SecKey)

}

func WatchSecProductInfoKey(key string)  {
	key1WatchChan := CliEtcd.Watch(context.Background(), key)
	var secProductInfo []model.SecInfo
	var getConfSucc = true

	for wresp := range key1WatchChan {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

			if ev.Type == clientv3.EventTypeDelete {
				log.Warnf("key[%s]的config deleted", key)
				continue
			}

			if ev.Type == clientv3.EventTypePut && string(ev.Kv.Key) == key {
				err := json.Unmarshal(ev.Kv.Value, &secProductInfo)
				if err != nil {
					log.Errorf("key [%s] unmarshal[%v] failed, err is:%v", key, secProductInfo, err)
					getConfSucc = false
					continue
				}
			}
			log.Debugf("get config from etcd, type:%v, key:%v, value:%v", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}

		if getConfSucc {
			log.Debugf("get config from etcd success, %v", secProductInfo)
			updateSecProductInfo(secProductInfo)
		}
	}
}

func updateSecProductInfo(infos []model.SecInfo)  {
	tmp := make(map[int] *model.SecInfo, 1024)
	for _, v := range infos {
		tmp[v.ProductId] = &v
	}
	RWLockerOfSecInfo.Lock()
	SecInfosMap = tmp
	RWLockerOfSecInfo.Unlock()
}

