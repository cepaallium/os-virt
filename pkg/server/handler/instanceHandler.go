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

type InstanceHandler struct {
	kubernetes.Client
}

func NewInstanceHandler() *InstanceHandler {
	c := new(InstanceHandler)
	c.Client = clients.Interface().Kubernetes(constants.ManagerCluster)
	return  c
}

func (h *InstanceHandler) GetInstanceDetail(param *param.InstanceParam) (*res.InstanceDetail, *errcode.ErrorInfo) {
	detail := &res.InstanceDetail{}
	detail.ClusterId = param.ClusterId
	detail.Uid = uuid.New().String()
	detail.InstanceId = param.InstanceId
	detail.InstanceName = param.InstanceName
	detail.Status = "OK"
	return detail, nil
}
