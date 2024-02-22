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

type TemplateHandler struct {
	kubernetes.Client
}

func NewTemplateHandler() *TemplateHandler {
	c := new(TemplateHandler)
	c.Client = clients.Interface().Kubernetes(constants.ManagerCluster)
	return  c
}

func (h *TemplateHandler) GetTemplateDetail(param *param.TemplateParam) (*res.TemplateDetail, *errcode.ErrorInfo) {
	detail := &res.TemplateDetail{}
	detail.ClusterId = param.ClusterId
	detail.Uid = uuid.New().String()
	detail.TemplateId = param.TemplateId
	detail.TemplateName = param.TemplateName
	detail.Status = "OK"
	return detail, nil
}
