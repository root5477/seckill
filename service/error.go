package service

import "errors"

const (
	ErrInvalidRequest    = 1001
	ErrNotFoundProductId = 1002
	ErrCheckAuthFaild    = 1003
	ErrUserServiceBusy   = 1004 //实际上是用户请求太快，返回这个
	ErrActiveNotStart    = 1005
	ErrActiveAlreadyEnd  = 1006
	ErrActiveSaleOut     = 1007
	ErrPleaseRetry       = 1008 //req处理超时
	ErrClientClosed      = 1009
)

const (
	ErrServiceBusy       = 2001
	ErrSecKillSucc       = 2002
	ErrNotFoundProductIdLayer = 2003
	ErrSoldOut           = 2004
	ErrRetry             = 2005
	ErrAlreadyBuyLimit   = 2006
)

var errMsg = map[int]string{
	ErrServiceBusy:     "服务器错误",
	ErrSecKillSucc:     "抢购成功",
	ErrNotFoundProductIdLayer: "没有该商品",
	ErrSoldOut:         "商品售罄",
	ErrRetry:           "请重试",
	ErrAlreadyBuyLimit: "已经抢购",
}

func GetErrMsg(code int) error {
	return errors.New(errMsg[code])
}