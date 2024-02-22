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

type VolumeHandler struct {
	kubernetes.Client
}

func NewVolumeHandler() *VolumeHandler {
	c := new(VolumeHandler)
	c.Client = clients.Interface().Kubernetes(constants.ManagerCluster)
	return  c
}

func (h *VolumeHandler) GetVolumeDetail(param *param.VolumeParam) (*res.VolumeDetail, *errcode.ErrorInfo) {
	detail := &res.VolumeDetail{}
	detail.ClusterId = param.ClusterId
	detail.Uid = uuid.New().String()
	detail.VolumeId = param.VolumeId
	detail.VolumeName = param.VolumeName
	detail.Status = "OK"
	return detail, nil
}
