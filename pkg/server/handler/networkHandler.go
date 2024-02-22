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

type NetworkHandler struct {
	kubernetes.Client
}

func NewNetworkHandler() *NetworkHandler {
	c := new(NetworkHandler)
	c.Client = clients.Interface().Kubernetes(constants.ManagerCluster)
	return  c
}

func (h *NetworkHandler) GetNetworkDetail(param *param.NetworkParam) (*res.NetworkDetail, *errcode.ErrorInfo) {
	detail := &res.NetworkDetail{}
	detail.ClusterId = param.ClusterId
	detail.Uid = uuid.New().String()
	detail.NetworkId = param.NetworkId
	detail.NetworkName = param.NetworkName
	detail.Status = "OK"
	return detail, nil
}
