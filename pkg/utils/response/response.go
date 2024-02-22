package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os-virt/pkg/utils/errcode"
)

type RestfulEntity struct {
	Status   bool        `json:"status"`
	Code     int         `json:"code"`
	Data     interface{} `json:"data"`
	Msg      string      `json:"msg"`
	HttpCode int         `json:"http_code"`
	Target   string      `json:"target"`
}

func SuccessReturn(c *gin.Context, obj interface{}) {
	res := RestfulEntity{
		Msg:      "success",
		Status:   true,
		Code:     http.StatusOK,
		Data:     obj,
		HttpCode: http.StatusOK,
	}
	c.JSON(http.StatusOK, res)
	c.Abort()
}

func FailReturn(c *gin.Context, errorInfo *errcode.ErrorInfo) {
	res := RestfulEntity{
		Msg:      errorInfo.Message,
		Status:   false,
		Code:     errorInfo.Code,
		Data:     nil,
		HttpCode: errorInfo.HttpCode,
	}
	c.JSON(errorInfo.HttpCode, res)
	c.Abort()
}
