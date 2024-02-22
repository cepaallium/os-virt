package apis

import (
	"github.com/gin-gonic/gin"
	"os-virt/pkg/server/handler"
	"os-virt/pkg/server/param"
	"os-virt/pkg/utils/errcode"
	"os-virt/pkg/utils/response"
)

type TemplateHandler struct {
	handler *handler.TemplateHandler
}

func NewTemplateApi() *TemplateHandler {
	c := new(TemplateHandler)
	c.handler = handler.NewTemplateHandler()
	return  c
}

func (c TemplateHandler) GetTemplateDetail(context *gin.Context) {

	data, isValid := param.CheckTemplateParam(context)
	if !isValid {
		response.FailReturn(context,errcode.InvalidParamValue)
		return
	}

	result, errInfo := c.handler.GetTemplateDetail(data)
	if errInfo != nil{
		response.FailReturn(context, errInfo)
		return
	}
	response.SuccessReturn(context, result)
}
