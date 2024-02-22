package apis

import (
	"github.com/gin-gonic/gin"
	"os-virt/pkg/server/handler"
	"os-virt/pkg/server/param"
	"os-virt/pkg/utils/errcode"
	"os-virt/pkg/utils/response"
)

type ImageApi struct {
	handler *handler.ImageHandler
}

func NewImageApi() *ImageApi {
	c := new(ImageApi)
	c.handler = handler.NewImageHandler()
	return  c
}

func (c ImageApi) GetImageDetail(context *gin.Context) {

	data, isValid := param.CheckImageParam(context)
	if !isValid {
		response.FailReturn(context,errcode.InvalidParamValue)
		return
	}

	result, errInfo := c.handler.GetImageDetail(data)
	if errInfo != nil{
		response.FailReturn(context, errInfo)
		return
	}
	response.SuccessReturn(context, result)
}
