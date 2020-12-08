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
		status = "the sec kill is not start"
	}
	if now - v.StartTime > 0 {
		startTag = true
	}
	if now - v.EndTime > 0 {
		startTag = false
		endTag = true
		status = "the sec kill is already end"
	}

	if v.Status == model.ProductStatusForceSaleOut || v.Status == model.ProductStatusSaleOut {
		startTag = false
		endTag = true
		status = "product is sale out"
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

	return
}
