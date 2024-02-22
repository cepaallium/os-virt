package apis

import (
	"github.com/gin-gonic/gin"
	"os-virt/pkg/server/handler"
	"os-virt/pkg/server/param"
	"os-virt/pkg/utils/errcode"
	"os-virt/pkg/utils/response"
)

type InstanceHandler struct {
	handler *handler.InstanceHandler
}

func NewInstanceApi() *InstanceHandler {
	c := new(InstanceHandler)
	c.handler = handler.NewInstanceHandler()
	return  c
}

func (c InstanceHandler) GetInstanceDetail(context *gin.Context) {

	data, isValid := param.CheckInstanceParam(context)
	if !isValid {
		response.FailReturn(context,errcode.InvalidParamValue)
		return
	}

	result, errInfo := c.handler.GetInstanceDetail(data)
	if errInfo != nil{
		response.FailReturn(context, errInfo)
		return
	}
	response.SuccessReturn(context, result)
}
