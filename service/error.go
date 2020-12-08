package service

const (
	ErrInvalidRequest = 1001
	ErrNotFoundProductId = 1002
	ErrCheckAuthFaild = 1003
	ErrUserServiceBusy = 1004  //实际上是用户请求太快，返回这个
)
