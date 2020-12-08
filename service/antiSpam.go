package service

import (
	"fmt"
	"secProxy/conf"
	"secProxy/model"
	"sync"
)

var (
	secLimitMgr = &SecLimitMgr{
		UserLimitMap:make(map[string] *SecLimit),
		IpLimitMap:make(map[string] *SecLimit),
	}
)

type SecLimitMgr struct {
	UserLimitMap map[string] *SecLimit
	IpLimitMap map[string] *SecLimit
	Locker sync.Mutex
}

//访问次数 计数控制，有两个方法
type SecLimit struct {
	count int
	currentTime int64
}

func (p *SecLimit) Count(nowTime int64) (couCount int) {
	//统计1秒内的访问次数，如果不是同一秒，将访问次数count重置为1
	if p.currentTime != nowTime {
		p.count = 1
		p.currentTime = nowTime
		couCount = p.count
		return
	}
	p.count ++
	couCount = p.count
	return
}

func (p *SecLimit) CheckNumNow(nowTime int64) int {
	if nowTime != p.currentTime {
		return 0
	}
	return p.count
}

func AntiSpam(cookie *model.UserCookie, req *model.SecKillReq) (err error) {
	secLimitMgr.Locker.Lock()
	defer secLimitMgr.Locker.Unlock()

	//1.针对用户身份的反作弊
	userSecLimit, ok := secLimitMgr.UserLimitMap[cookie.UserId]
	if !ok {
		userSecLimit = &SecLimit{}
		secLimitMgr.UserLimitMap[cookie.UserId] = userSecLimit
	}
	count := userSecLimit.Count(req.AccessTime)
	if count > conf.KillConf.Limit.UserSecAccessLimit {
		err = fmt.Errorf("invalid request1")
		return
	}

	//2.针对客户端ip的反作弊
	ipSecLimit, ok := secLimitMgr.IpLimitMap[req.ClientAddr]
	if !ok {
		ipSecLimit = &SecLimit{}
		secLimitMgr.IpLimitMap[req.ClientAddr] = ipSecLimit
	}
	count = ipSecLimit.Count(req.AccessTime)
	if count > conf.KillConf.Limit.IpSecAccessLimit {
		err = fmt.Errorf("invalid request2")
		return
	}
	return
}
