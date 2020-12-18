package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"secProxy/model"
	"testing"
	"time"
)

func TestEnvInit(t *testing.T) {
	EnvInit("/Users/chenqi/go/src/secProxy/conf/seckill.json")
	fmt.Println(*KillConf)
	fmt.Println("secKey:", KillConf.Etcd.SecKey)
}

func TestInitRedis(t *testing.T) {
	err := InitRedis()
	if err != nil {
		t.Errorf("test failed!")
		return
	}
	t.Logf("test success!")
}

func TestInitEtcd(t *testing.T) {
	err := InitEtcd()
	if err != nil {
		t.Errorf("TestInitEtcd failed, err is:%v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	_, err = CliEtcd.Put(ctx, "key1", "value--1")
	cancel()
	if err != nil {
		t.Errorf("TestInitEtcd failed, err2 is:%v", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp ,err := CliEtcd.Get(ctx, "key1")
	cancel()
	if err != nil {
		t.Errorf("TestInitEtcd failed, err3 is:%v", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
}

func TestWatchEtcdKey(t *testing.T)  {
	err := InitEtcd()
	if err != nil {
		panic(err)
	}
	key1WatchChan := CliEtcd.Watch(context.Background(), "key1")
	for wresp := range key1WatchChan {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func TestSetSecInfo(t *testing.T) {

	err := InitEtcd()
	if err != nil {
		panic(err)
	}
	product1 := &model.SecInfo{
		ProductId:1,
		StartTime:time.Now().Unix(),
		EndTime:time.Now().Unix() + 60 * 60 * 24,
		Status:0,
		Total:10000,
		Left:10000,
	}
	product2 := &model.SecInfo{
		ProductId:2,
		StartTime:time.Now().Unix(),
		EndTime:time.Now().Unix() + 60 * 60 * 24,
		Status:0,
		Total:10000,
		Left:10000,
	}
	var products []*model.SecInfo

	products = append(products, product1, product2)
	bytes, err := json.Marshal(&products)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	fmt.Println("23333:", KillConf.Etcd.SecKey)
	fmt.Println( string(bytes))
	_, err = CliEtcd.Put(ctx, KillConf.Etcd.SecKey, string(bytes))
	cancel()
	if err != nil {
		panic(err)
	}
}

func InitEtcd222() (CliEtcd *clientv3.Client, err error) {

	CliEtcd, err = clientv3.New(clientv3.Config{
		Endpoints:[]string {"10.226.133.69:2379"},
		DialTimeout: time.Duration(KillConf.Etcd.TimeOut) * time.Second,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err is:%v\n", err)
		return nil, err
	}
	fmt.Println("connect to etcd success!")
	return
}


func TestSetSecInfo2(t *testing.T) {

	CliEtcd2,err := InitEtcd222()
	if err != nil {
		panic(err)
	}
	product1 := &model.SecInfo{
		ProductId:1,
		StartTime:time.Now().Unix(),
		EndTime:time.Now().Unix() + 60 * 60 * 24,
		Status:0,
		Total:10000,
		Left:10000,
	}
	product2 := &model.SecInfo {
		ProductId:2,
		StartTime:time.Now().Unix(),
		EndTime:time.Now().Unix() + 60 * 60 * 24,
		Status:0,
		Total:10000,
		Left:10000,
	}
	var products []*model.SecInfo

	products = append(products, product1, product2)
	bytes, err := json.Marshal(&products)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	fmt.Println( string(bytes))
	_, err = CliEtcd2.Put(ctx, "sec_product_info", string(bytes))
	defer cancel()
	if err != nil {
		panic(err)
	}
	fmt.Println("hahahahahhahahah")
}


