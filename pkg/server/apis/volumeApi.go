package apis

import (
	"github.com/gin-gonic/gin"
	"os-virt/pkg/server/handler"
	"os-virt/pkg/server/param"
	"os-virt/pkg/utils/errcode"
	"os-virt/pkg/utils/response"
)

type VolumeApi struct {
	handler *handler.VolumeHandler
}

func NewVolumeApi() *VolumeApi {
	c := new(VolumeApi)
	c.handler = handler.NewVolumeHandler()
	return  c
}

func (c VolumeApi) GetVolumeDetail(context *gin.Context) {

	data, isValid := param.CheckVolumeParam(context)
	if !isValid {
		response.FailReturn(context,errcode.InvalidParamValue)
		return
	}

	result, errInfo := c.handler.GetVolumeDetail(data)
	if errInfo != nil{
		response.FailReturn(context, errInfo)
		return
	}
	response.SuccessReturn(context, result)
}
