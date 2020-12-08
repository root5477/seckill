package model

type SecKillReq struct {
	Product     int    `json:"product"`
	Src         string `json:"src"`
	AuthCode    string `json:"auth_code"`
	Time        string `json:"time"`
	Nance       string `json:"nance"`
	AccessTime  int64  `json:"access_time"`
	ClientAddr  string `json:"client_addr"`
	ClientRefer string `json:"client_refer"`
}

type UserCookie struct {
	UserId         string `json:"user_id"`
	UserCookieAuth string `json:"user_cookie_auth"`
}
