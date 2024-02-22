/*
Copyright 2021 KubeCube Authors

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

package clients

import (
	"kubevirt.io/client-go/kubecli"
	"os-virt/pkg/clients/kubernetes"
	"os-virt/pkg/multicluster"
	"os-virt/pkg/multicluster/manager"
	"os-virt/pkg/utils/log"
)

// Clients aggregates all clients of cube needed
type Clients interface {
	Kubernetes(cluster string) kubernetes.Client
}

// genericClientSet is global cube cube client that must init at first
var genericClientSet = &cceClientSet{}

type cceClientSet struct {
	k8s manager.ClustersManagerInterface
}

// InitCceManagerClientSetWithOpts initialize global clients with given config.
func InitCceManagerClientSetWithOpts() {
	genericClientSet.k8s = multicluster.Interface()
}

// Interface the entry for cube client
func Interface() *cceClientSet {
	return genericClientSet
}

// Kubernetes get the indicate client for cluster, we log error instead of return it
// for convenience, caller needs to determine whether the return value is nil
func (c *cceClientSet) Kubernetes(clusterId string) kubernetes.Client {
	client, err := c.k8s.GetClient(clusterId)
	if err != nil {
		log.Error("get cluster client error,clusterId: %s", clusterId)
		return nil
	}

	return client
}
func (c *cceClientSet) KubeVirt(clusterId string) kubecli.KubevirtClient {
	client, err := c.k8s.GetKubeVirtClient(clusterId)
	if err != nil {
		log.Error("get cluster kubeVirt client error, clusterId: %s", clusterId)
		return nil
	}
	return client
}
