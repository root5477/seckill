package log

import (
	"secProxy/conf"
	"fmt"
)

func InitLog(path string, conf *conf.Config) error {
	err := Init(path, conf.Log.MinLevel)
	if err != nil {
		return fmt.Errorf("Log init error:%v", err)
	}

	Info("Log init success")
	return nil
}
