package handler

import (
	"github.com/google/uuid"
	"os-virt/pkg/clients"
	"os-virt/pkg/clients/kubernetes"
	"os-virt/pkg/server/param"
	"os-virt/pkg/server/res"
	"os-virt/pkg/utils/constants"
	"os-virt/pkg/utils/errcode"
)

type ImageHandler struct {
	kubernetes.Client
}

func NewImageHandler() *ImageHandler {
	c := new(ImageHandler)
	c.Client = clients.Interface().Kubernetes(constants.ManagerCluster)
	return  c
}

func (h *ImageHandler) GetImageDetail(param *param.ImageParam) (*res.ImageDetail, *errcode.ErrorInfo) {
	detail := &res.ImageDetail{}
	detail.ClusterId = param.ClusterId
	detail.Uid = uuid.New().String()
	detail.ImageId = param.ImageId
	detail.ImageName = param.ImageName
	detail.Status = "OK"
	return detail, nil
}
