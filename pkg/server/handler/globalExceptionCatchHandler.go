package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"os-virt/pkg/utils/errcode"
	"os-virt/pkg/utils/response"
	"runtime/debug"
)

func GlobalExceptionCatch(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			msg, _ := json.Marshal(err)
			log.Printf("出现异常： %s\n%s", msg, debug.Stack())
			var ex *errcode.ErrorInfo
			if h, ok := err.(*errcode.ErrorInfo); ok {
				ex = h
			} else {
				ex = errcode.CustomReturn(errcode.InternalServerError.HttpCode, string(msg))
			}
			response.FailReturn(c, ex)
		}
	}()
	c.Next()
}
