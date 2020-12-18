package model

type SecKillReq struct {
	Product     int       `json:"product"`
	Src         string    `json:"src"`
	AuthCode    string    `json:"auth_code"`
	Time        string    `json:"time"`
	Nance       string    `json:"nance"`
	AccessTime  int64     `json:"access_time"`
	ClientAddr  string    `json:"client_addr"`
	ClientRefer string    `json:"client_refer"`
	CloseNotify <- chan bool `json:"close_notify"`
}

type UserCookie struct {
	UserId         string `json:"user_id"`
	UserCookieAuth string `json:"user_cookie_auth"`
}

type SecRequestWithCookie struct {
	SecKillReq
	UserCookie
}

type SecResponse struct {
	ProductId int    `json:"product_id"`
	UserId    string `json:"user_id"`
	Token     string `json:"token"`
	TokenTime int64  `json:"token_time"`
	Code      int    `json:"code"` //用于标示是否抢到
}
