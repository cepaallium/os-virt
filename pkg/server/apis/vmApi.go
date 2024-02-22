package apis

import (
	"github.com/gin-gonic/gin"
	"kubevirt.io/api/core/v1"
	"os-virt/pkg/server/handler"
	"os-virt/pkg/server/param"
	"os-virt/pkg/utils/constants"
	"os-virt/pkg/utils/errcode"
	"os-virt/pkg/utils/log"
	"os-virt/pkg/utils/page"
	"os-virt/pkg/utils/response"
	"sort"
)

type VmApi struct {
	vmHandler *handler.VmHandler
}

func NewVmApi() *VmApi {
	c := new(VmApi)
	c.vmHandler = handler.NewVmHandler()
	return c
}

func (c VmApi) GetVmDetail(context *gin.Context) {
	vmName := context.Param("vmName")
	clusterId := context.Query("clusterId")
	namespace := context.Query("namespace")
	if len(clusterId) == 0 || len(vmName) == 0 || len(namespace) == 0 {
		panic(errcode.InvalidParamValue)
	}

	result := c.vmHandler.GetVmDetail(clusterId, namespace, vmName)
	response.SuccessReturn(context, result)
}

func (c VmApi) CreateVm(context *gin.Context) {
	vmParam := &param.CreateVmParam{}
	err := context.ShouldBindJSON(vmParam)
	if err != nil {
		log.Error("data bind error.", err)
		panic(errcode.InvalidParamValue)
	}

	c.vmHandler.CreateVm(vmParam)
	response.SuccessReturn(context, vmParam.Name)
}

func (c VmApi) GetVms(context *gin.Context) {
	vms := c.vmHandler.GetVms(context)
	sortVms := sortVms(context, vms)
	doPage := page.DoPage(context.Query(constants.CurrentPage), context.Query(constants.PageSize), sortVms)
	response.SuccessReturn(context, doPage)
}

func (c VmApi) DeleteVm(context *gin.Context) {
	vmName := context.Param("vmName")
	clusterId := context.Query("clusterId")
	namespace := context.Query("namespace")
	c.vmHandler.DeleteVm(clusterId, namespace, vmName)
	response.SuccessReturn(context, vmName)
}

func sortVms(context *gin.Context, vms []v1.VirtualMachine) []interface{} {
	order := context.Query("order")
	sort.Slice(vms, func(i, j int) bool {
		result := vms[i].CreationTimestamp.Time.Unix() < vms[j].CreationTimestamp.Time.Unix()
		if vms[i].CreationTimestamp.Time.Unix() == vms[j].CreationTimestamp.Time.Unix() {
			result = vms[i].Name < vms[j].Name
		}
		if order == "" {
			return result
		}
		if order == "desc" {
			return !result
		}
		return result
	})
	var records []interface{}
	for _, item := range vms {
		records = append(records, item)
	}
	return records
}
