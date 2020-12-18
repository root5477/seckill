package service

import (
	"crypto/md5"
	"fmt"
	"secProxy/conf"
	"secProxy/log"
	"secProxy/model"
	"time"
)

var (
	SecKillConf *model.SecInfo
)

func InitService(serviceConf *model.SecInfo)  {

}

func GetSecInfoById(productId int) (item map[string] interface{}, code int, err error) {
	v, ok := conf.SecInfosMap[productId]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found productId:%v", productId)
		return
	}
	now := time.Now().Unix()
	startTag := false
	endTag := false
	status := "success"

	item = make(map[string] interface{})
	item["product_id"] = v.ProductId

	if now - v.StartTime < 0 {
		startTag = false
		endTag = false
		status = "the sec kill is not start"
		code = ErrActiveNotStart
	}
	if now - v.StartTime > 0 {
		startTag = true
	}
	if now - v.EndTime > 0 {
		startTag = false
		endTag = true
		status = "the sec kill is already end"
		code = ErrActiveAlreadyEnd
	}

	if v.Status == model.ProductStatusForceSaleOut || v.Status == model.ProductStatusSaleOut {
		startTag = false
		endTag = true
		status = "product is sale out"
		code = ErrActiveSaleOut
	}

	item["start"] = startTag
	item["end"] = endTag
	item["status"] = status
	return
}

func SecInfoList() (data []map[string]interface{}) {
	//加读锁
	conf.RWLockerOfSecInfo.RLock()
	defer conf.RWLockerOfSecInfo.RUnlock()
	for _, v := range conf.SecInfosMap {
		info, _, err := GetSecInfoById(v.ProductId)
		if err != nil {
			log.Errorf("get product_id[%v] failed, err is:%v", v.ProductId, err)
			continue
		}
		data = append(data, info)
	}
	return
}

func SecInfo(productId int) (data []map[string]interface{}, code int, err error){
	//加读锁
	conf.RWLockerOfSecInfo.RLock()
	defer conf.RWLockerOfSecInfo.RUnlock()
	info, code, err := GetSecInfoById(productId)
	if err != nil {
		return
	}
	data = append(data, info)
	return
}

func UserCheck(cookie *model.UserCookie, req *model.SecKillReq) (err error) {
	//1.校验refer是否在白名单
	found := false
	for _, refer := range conf.KillConf.Limit.ReferWhiteList {
		if refer == req.ClientRefer {
			found = true
			break
		}
	}
	if !found {
		log.Warnf("refer [%v] not in whitelist, req:%v, cookie:%v", req.ClientRefer, req, cookie)
		err = fmt.Errorf("invalid request")
		return
	}


	authData := fmt.Sprintf("%v-%v", cookie.UserId, conf.KillConf.SecretKey.CookieSecretKey)
	authSign := fmt.Sprintf("%x", md5.Sum([]byte(authData)))
	if authSign != cookie.UserCookieAuth {
		err = fmt.Errorf("invalid user cookie auth")
		return
	}
	return
}

func SecKill(req *model.SecKillReq, cookie *model.UserCookie) (data interface{}, code int, err error) {
	conf.RWLockerOfSecInfo.Lock()
	defer conf.RWLockerOfSecInfo.Unlock()

	//black ip&&id check
	_, ok := conf.KillConf.IpBlackMap[req.ClientAddr]
	if ok {
		log.Errorf("req client ip[%v] is black, req:[%v]", req.ClientAddr, *req)
		code = ErrCheckAuthFaild
		err = fmt.Errorf("invalid request")
		return
	}
	_, ok = conf.KillConf.IdBlackMap[cookie.UserId]
	if ok {
		log.Errorf("req user id is black, cookie:[%v], req:[%v]", *cookie, *req)
		code = ErrCheckAuthFaild
		err = fmt.Errorf("invalid request")
		return
	}

	err = UserCheck(cookie, req)
	if err != nil {
		log.Errorf("UserCheck failed, userId[%v] invalid, req:[%v]", cookie.UserId, req)
		code = ErrCheckAuthFaild
		return
	}

	err = AntiSpam(cookie, req)
	if err != nil {
		code = ErrUserServiceBusy
		log.Errorf("userId[%v] request too many in one second, req:[%v]", cookie.UserId, req)
		return
	}

	data, code, err =  GetSecInfoById(req.Product)
	if err != nil {
		log.Errorf("userId[%v] GetSecInfoById failed, req:[%v]", *cookie, *req)
		return
	}

	if code != 0 {
		log.Warnf("userId[%v] GetSecInfoById failed, req:[%v]", *cookie, *req)
		return
	}

	//到这里说明请求合法 && 活动正常进行中, 发送给redis队列处理
	reqWithCookie := &model.SecRequestWithCookie{
		SecKillReq:*req,
		UserCookie:*cookie,
	}
	conf.SecReqChan <- reqWithCookie
	return
}

func InitRedisProcessFunc()  {
	for i := 0; i < conf.KillConf.Server.WriteProxy2LayerGoroutineNum; i ++ {
		go WriteHandle()
	}

	for i := 0; i < conf.KillConf.Server.ReadLayer2ProxyGoroutineNum; i ++ {
		go ReadHandle()
	}
}


