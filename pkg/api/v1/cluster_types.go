/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"os-virt/pkg/utils/kubeconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ClusterState string

const (
	// ClusterInitFailed happened when init cluster failed
	// generally when network error occurred
	ClusterInitFailed ClusterState = "initFailed"

	// ClusterProcessing wait for cluster be taken over
	ClusterProcessing ClusterState = "processing"

	// ClusterDeleting means cluster is under deleting
	ClusterDeleting ClusterState = "deleting"

	// ClusterNormal represent cluster is healthy
	ClusterNormal ClusterState = "normal"

	// ClusterUnormal represent cluster is unhealthy
	ClusterUnormal ClusterState = "unhealthy"
)

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	ClusterId string `json:"clusterId,omitempty"`
	// KubeConfig contains cluster raw kubeConfig
	KubeConfig []byte `json:"kubeconfig,omitempty"`

	// Kubernetes API Server endpoint. Example: https://10.10.0.1:6443
	KubernetesAPIEndpoint string `json:"kubernetesAPIEndpoint,omitempty"`

	// cluster type  member or manager
	ClusterType string `json:"clusterType,omitempty"`

	// describe cluster
	// +optional
	Description string `json:"description,omitempty"`
}

type WorkNode struct {
	IpAddr string `json:"ipAddr"`
	Status string `json:"status"`
	CpuUsedPercent string `json:"cpuUsedPercent"`
	MemoryUsedPercent string `json:"memoryUsedPercent"`
	DiskUsedPercent string `json:"diskUsedPercent"`
	PodCount int `json:"podCount"`
	VolumeCount int `json:"volumeCount"`
}

type MasterNode struct {
	IpAddr string `json:"ipAddr"`
	Status string `json:"status"`
	CpuUsedPercent string `json:"cpuUsedPercent"`
	MemoryUsedPercent string `json:"memoryUsedPercent"`
	DiskUsedPercent string `json:"diskUsedPercent"`
	PodCount int `json:"podCount"`
	VolumeCount int `json:"volumeCount"`
	IsVip bool `json:"isVip"`
}

type ClusterInfo struct {
	ClusterId string `json:"clusterId,omitempty"`
	ClusterName string         `json:"clusterName,omitempty"`
	State *ClusterState        `json:"clusterState,omitempty"`
	MasterList []MasterNode    `json:"masterList,omitempty"`
	NodeList []WorkNode        `json:"nodeList,omitempty"`
	LastHeartbeat *metav1.Time `json:"lastHeartbeat,omitempty"`
}


// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	ClusterInfo *ClusterInfo `json:"clusterInfo,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

func (c *Cluster) GenerateClient()(client.Client,error) {
	config, err := kubeconfig.LoadKubeConfigFromBytes(c.Spec.KubeConfig)
	if err != nil {
		return nil,err
	}
	scheme := runtime.NewScheme()
	utilruntime.Must(AddToScheme(scheme))
	cli, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}
	return  cli,nil
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
