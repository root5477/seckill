package dao

import (
	"fmt"
	"secProxy/conf"
	"testing"
)

func TestLoadSecInfoConf(t *testing.T) {
	err := LoadSecInfoConf()
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.SecInfosMap)
	fmt.Println(conf.SecInfosMap[1])
	fmt.Println(conf.SecInfosMap[2])
}

func TestTmp(t *testing.T)  {
	type ins struct {
		Id int
		Value string
	}
	ins1 := ins {
		Id:1,
		Value:"va1",
	}
	ins2 := ins{
		Id:2,
		Value:"va2",
	}
	ins3 := ins{
		Id:3,
		Value:"va3",
	}
	var infoSlice []ins
	infoSlice = append(infoSlice, ins1, ins2, ins3)

	var m = make(map[int] *ins)
	for _, v := range infoSlice {
		m[v.Id] = &v
	}
	fmt.Println(m[1], m[2], m[3])
}