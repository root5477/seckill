package log

import (
	"secProxy/conf"
	"testing"
)

func TestInitLog(t *testing.T) {
	err := InitLog("./seelog.xml", conf.KillConf)
	if err != nil {
		panic(err)
	}
	defer Flush()
	Debug("test success !")
}
