package param

import (
	"github.com/gin-gonic/gin"
	"os-virt/pkg/utils/log"
)

func CheckImageParam(context *gin.Context) (*ImageParam, bool) {
	param := &ImageParam{}
	err := context.ShouldBindJSON(param)
	if err != nil {
		log.Error(err.Error())
		return nil, false
	}
	return param, true
}

func CheckInstanceParam(context *gin.Context) (*InstanceParam, bool) {
	param := &InstanceParam{}
	err := context.ShouldBindJSON(param)
	if err != nil {
		log.Error(err.Error())
		return nil, false
	}
	return param, true
}

func CheckNetworkParam(context *gin.Context) (*NetworkParam, bool) {
	param := &NetworkParam{}
	err := context.ShouldBindJSON(param)
	if err != nil {
		log.Error(err.Error())
		return nil, false
	}
	return param, true
}

func CheckTemplateParam(context *gin.Context) (*TemplateParam, bool) {
	param := &TemplateParam{}
	err := context.ShouldBindJSON(param)
	if err != nil {
		log.Error(err.Error())
		return nil, false
	}
	return param, true
}

func CheckVolumeParam(context *gin.Context) (*VolumeParam, bool) {
	param := &VolumeParam{}
	err := context.ShouldBindJSON(param)
	if err != nil {
		log.Error(err.Error())
		return nil, false
	}
	return param, true
}
