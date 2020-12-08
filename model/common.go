package model

type CommonResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const  (
	ProductStatusNormal = 0
	ProductStatusSaleOut = 1
	ProductStatusForceSaleOut = 2
)