package handler

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubevirt.io/api/core/v1"
	"os"
	"os-virt/pkg/clients"
	"os-virt/pkg/clients/kubernetes"
	"os-virt/pkg/server/param"
	"os-virt/pkg/server/res"
	"os-virt/pkg/utils/errcode"
	"os-virt/pkg/utils/log"
	"strings"
)

type VmHandler struct {
	K8SClient kubernetes.Client
	//KubeVirtClient kubecli.KubevirtClient
}

func NewVmHandler() *VmHandler {
	c := new(VmHandler)
	return c
}

func (h *VmHandler) GetVmDetail(clusterId, namespace, vmName string) *res.VmDetail {

	kubeVirtClient := clients.Interface().KubeVirt(clusterId)

	virtualMachine, err := kubeVirtClient.VirtualMachine(namespace).Get(context.Background(), vmName, &metav1.GetOptions{})

	if err != nil {
		log.Error("get namespace error: %v", err)
		if errors.IsNotFound(err) {
			panic(errcode.ResourceNotFound)
		}
		panic(err)
	}

	vmDetail := &res.VmDetail{}
	vmDetail.ClusterId = clusterId
	vmDetail.VmName = virtualMachine.Name
	vmDetail.Status = string(virtualMachine.Status.PrintableStatus)
	vmDetail.Uid = string(virtualMachine.UID)
	vmDetail.CreationTime = virtualMachine.GetObjectMeta().GetCreationTimestamp().String()

	return vmDetail
}

func (h *VmHandler) CreateVm(param *param.CreateVmParam) {

	clusterId := param.ClusterId
	kubeVirtClient := clients.Interface().KubeVirt(clusterId)

	resourceList := corev1.ResourceList{}
	resourceList[corev1.ResourceCPU] = resource.MustParse(param.Cpu)
	resourceList[corev1.ResourceMemory] = resource.MustParse(param.Memory)

	vmRunningTrue := true
	vm := v1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: param.Namespace,
			Name:      param.Name,
		},
		Spec: v1.VirtualMachineSpec{
			Running: &vmRunningTrue,
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				Spec: v1.VirtualMachineInstanceSpec{
					Domain: v1.DomainSpec{
						Resources: v1.ResourceRequirements{
							Requests: resourceList,
						},
						Devices: v1.Devices{
							Disks: []v1.Disk{
								{
									Name: "cloudinitdisk",
									DiskDevice: v1.DiskDevice{
										Disk: &v1.DiskTarget{},
									},
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "cloudinitdisk",
							VolumeSource: v1.VolumeSource{
								CloudInitNoCloud: &v1.CloudInitNoCloudSource{
									UserDataBase64: buildAdduserStartupScript(param.LoginConfig.Username, param.LoginConfig.Password),
								},
							},
						},
					},
					Networks: []v1.Network{
						{},
					},
				},
			},
		},
	}

	_, err := kubeVirtClient.VirtualMachine(param.Namespace).Create(context.Background(), &vm)

	if err != nil {
		log.Error("create statefulSet err, %v", err)

		if errors.IsAlreadyExists(err) {
			log.Error(" statefulSet %s in namespace %s already exist.%v", param.Name, param.Namespace, err)
			panic(errcode.ResourceAlreadyExists)
		} else {
			log.Error("create statefulSet %s in namespace %s error.%v", param.Name, param.Namespace, err)
			panic(err)
		}
	}
}

func buildAdduserStartupScript(username, password string) string {
	adduserStartupScriptByte, err := os.ReadFile("config/startup-script-adduser.sh")
	if err != nil {
		log.Error("read config/startup-script-adduser.sh error.", err)
		panic(err)
	}

	adduserStartupScript := string(adduserStartupScriptByte)

	strings.ReplaceAll(adduserStartupScript, "#{NEW_USER}", username)
	strings.ReplaceAll(adduserStartupScript, "#{NEW_USER_PASSWD}", password)

	adduserStartupScriptBase64 := base64.StdEncoding.EncodeToString([]byte(adduserStartupScript))
	return adduserStartupScriptBase64
}

func (h *VmHandler) GetVms(ctx *gin.Context) []v1.VirtualMachine {

	clusterId := ctx.Query("clusterId")
	namespace := ctx.Query("namespace")
	name := ctx.Query("name")

	kubeVirtClient := clients.Interface().KubeVirt(clusterId)

	virtualMachineList, err := kubeVirtClient.VirtualMachine(namespace).List(context.Background(), &metav1.ListOptions{})

	if err != nil {
		log.Error("get virtualMachineList %s error,%v", name, err)
		panic(err)
	}

	machines := make([]v1.VirtualMachine, 0)

	for _, item := range virtualMachineList.Items {
		if len(name) > 0 && !strings.Contains(item.Name, name) {
			continue
		}
		machines = append(machines, item)
	}

	return machines
}

func (h *VmHandler) DeleteVm(clusterId, namespace, name string) {
	kubeVirtClient := clients.Interface().KubeVirt(clusterId)
	err := kubeVirtClient.VirtualMachine(namespace).Delete(context.Background(), name, &metav1.DeleteOptions{})
	if err != nil {
		log.Error("delete VirtualMachine %s in namespace %s failed, cause: %s", name, namespace, err)
	}
}
