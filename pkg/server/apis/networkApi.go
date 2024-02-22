package apis

import (
	"github.com/gin-gonic/gin"
	"os-virt/pkg/server/handler"
	"os-virt/pkg/server/param"
	"os-virt/pkg/utils/errcode"
	"os-virt/pkg/utils/response"
)

type NetworkHandler struct {
	handler *handler.NetworkHandler
}

func NewNetworkApi() *NetworkHandler {
	c := new(NetworkHandler)
	c.handler = handler.NewNetworkHandler()
	return  c
}

func (c NetworkHandler) GetNetworkDetail(context *gin.Context) {

	data, isValid := param.CheckNetworkParam(context)
	if !isValid {
		response.FailReturn(context,errcode.InvalidParamValue)
		return
	}

	result, errInfo := c.handler.GetNetworkDetail(data)
	if errInfo != nil{
		response.FailReturn(context, errInfo)
		return
	}
	response.SuccessReturn(context, result)
}
