package app

import (
	"github.com/gin-gonic/gin"

	"net/http"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseSuccess(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, response{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	})
	return
}

// 系统未知错误http code 500
func ResponseFailed(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusInternalServerError, response{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Data:    data,
	})
	return
}

// 校验参数未通过http code 400
func ResponseCheckParamsFailed(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusBadRequest, response{
		Code:    http.StatusBadRequest,
		Message: msg,
		Data:    data,
	})
	return
}

// 未登录http code 401
func ResponseAuthFailed(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusUnauthorized, response{
		Code:    http.StatusUnauthorized,
		Message: msg,
		Data:    data,
	})
	return
}

// 未找到资源http code 404
func ResponseNoFound(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusNotFound, response{
		Code:    http.StatusNotFound,
		Message: msg,
		Data:    data,
	})
	return
}

// 未受权限http code 403
func ResponseNoPermission(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusForbidden, response{
		Code:    http.StatusForbidden,
		Message: msg,
		Data:    data,
	})
	return
}
