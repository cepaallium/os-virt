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

package manager

import (
	"kubevirt.io/client-go/kubecli"
	"os-virt/pkg/clients/kubernetes"
)

func (m *ClustersManager) GetClient(clusterId string) (kubernetes.Client, error) {
	c, err := m.Get(clusterId)
	if err != nil {
		return nil, err
	}

	return c.Client, err
}

func (m *ClustersManager) GetKubeVirtClient(clusterId string) (kubecli.KubevirtClient, error) {
	c, err := m.Get(clusterId)
	if err != nil {
		return nil, err
	}

	return c.VirtClient, err
}
