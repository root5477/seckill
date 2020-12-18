package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"secProxy/log"
	"secProxy/model"
	"secProxy/service"
	"strconv"
	"strings"
	"time"
)


//
func SecInfo(c *gin.Context)  {
	productIdStr := c.Query("product_id")
	productId, err := strconv.Atoi(productIdStr)
	resp := &model.CommonResp{}
	if err != nil {
		data := service.SecInfoList()
		resp.Code = 1001
		resp.Message = "success"
		resp.Data = data
		c.JSON(http.StatusBadRequest, resp)
		return
	} else {
		log.Debugf("get productId:%v", productId)
		data, code, err := service.SecInfo(productId)
		resp.Data = data
		resp.Code = code
		if err != nil {
			log.Errorf("get secInfo failed, err is:%v", err)
			resp.Message = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		resp.Message = "success"
		c.JSON(http.StatusOK, resp)
	}

}

func SecKillHandle(c *gin.Context)  {
	req := &model.SecKillReq{}
	cookie := &model.UserCookie{}
	resp := model.CommonResp{}

	req.Src = c.Query("src")
	req.Time = c.Query("time")
	req.AuthCode = c.Query("auth_code")
	req.Nance = c.Query("nance")


	productIdStr := c.Query("product")
	if productIdStr == "" || req.Src == "" || req.Nance == "" || req.Time == "" {
		log.Errorf("get params failed, req:%v", req)
		resp.Code = service.ErrInvalidRequest
		resp.Message = "req invalid ''"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	req.AccessTime = time.Now().Unix()
	//根据ip && refer进行安全策略判断
	if c.ClientIP() != "" {
		req.ClientAddr = strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	req.ClientRefer = c.Request.Referer()
	log.Debugf("client request:[%v]", req)
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		log.Errorf("get productId failed, err is:%v", err)
		resp.Code = service.ErrInvalidRequest
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	req.Product = productId

	userAuthSign, err1 := c.Cookie("userAuthSign")
	userId, err2 := c.Cookie("userId")
	if err1 != nil || err2 != nil {
		log.Errorf("get cookie failed, err1 is:[%v], err2 is:[%v]", err1, err2)
		resp.Code = service.ErrInvalidRequest
		resp.Message = "cookie is not right"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	cookie.UserCookieAuth = userAuthSign
	cookie.UserId = userId
	req.CloseNotify = c.Writer.CloseNotify()
	data, code, err := service.SecKill(req, cookie)
	if err != nil {
		resp.Code = code
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	log.Debugf("req:[%v], cookie:[%v], data:[%v]", req, cookie, data)

	resp.Code = 0
	resp.Message = "success"
	c.JSON(http.StatusOK, resp)
}