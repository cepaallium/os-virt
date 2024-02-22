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
	"context"
	"fmt"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"kubevirt.io/client-go/kubecli"
	ccev1 "os-virt/pkg/api/v1"
	"os-virt/pkg/clients/kubernetes"
	"os-virt/pkg/utils/constants"
	"os-virt/pkg/utils/kubeconfig"
	"os-virt/pkg/utils/log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sync"
)

// ClustersManagerInterface access to internal cluster
type ClustersManagerInterface interface {
	// Add runtime cache in memory
	Add(cluster string, internalCluster *InternalCluster) error
	Get(cluster string) (*InternalCluster, error)
	Del(cluster string) error

	// MonitorFor collector heartbeat for collector
	//MonitorFor(ctx context.Context, clusterId string) error

	AddInternalCluster(cluster *ccev1.Cluster) error
	// GetClient get client for cluster
	GetClient(clusterId string) (kubernetes.Client, error)

	GetKubeVirtClient(clusterId string) (kubecli.KubevirtClient, error)
}

// ClusterManagerInstance instance implement interface,
// init manager cluster at first.
var ClusterManagerInstance = newClusterMgr()

// newMultiClusterMgr init MultiClustersMgr with pivot internal cluster
func newClusterMgr() *ClustersManager {
	clustersManager := &ClustersManager{Clusters: make(map[string]*InternalCluster)}
	config, err := ctrl.GetConfig()
	if err != nil {
		log.Warn("get kubeconfig failed: %v", err)
		return nil
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	//utilruntime.Must(hnc.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
	utilruntime.Must(ccev1.AddToScheme(scheme))

	cli, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatal("connect to manager cluster failed: %v", err)
	}

	cluster := ccev1.Cluster{}
	err = cli.Get(context.Background(), types.NamespacedName{Name: constants.ManagerCluster, Namespace: "cloudos-cce-project"}, &cluster)
	if err != nil {
		log.Fatal("get manager cluster failed: %v", err)
	}

	cfg, err := kubeconfig.LoadKubeConfigFromBytes(cluster.Spec.KubeConfig)
	if err != nil {
		log.Fatal("invalid kubeconfig of manager cluster: %v", err)
	}

	c := new(InternalCluster)
	c.StopCh = make(chan struct{})
	c.Config = cfg
	c.Client, err = kubernetes.NewClientFor(cfg, c.StopCh)
	//c.Monitor = monitor.NewClusterMonitor(cluster.Spec.ClusterId, 0, 0, c.Client.Direct(), c.StopCh)
	if err != nil {
		// early exit when connect to the k8s apiserver of control plane failed
		log.Fatal("make client for manager cluster failed: %v", err)
	}
	err = clustersManager.Add(constants.ManagerCluster, c)
	if err != nil {
		log.Fatal("init multi cluster mgr failed: %v", err)
	}
	return clustersManager
}

// InternalCluster represent a cluster runtime contains
// client and internal collector.
type InternalCluster struct {
	// Client holds all the clients needed
	Client kubernetes.Client

	VirtClient kubecli.KubevirtClient

	// Config bind to a real cluster
	Config *rest.Config

	//Monitor *monitor.ClusterMonitor

	// StopCh for closing channel when delete cluster, goroutine
	// of informer and scout will exit gracefully.
	StopCh chan struct{}
}

// MultiClustersMgr a memory cache for runtime cluster.
type ClustersManager struct {
	sync.RWMutex
	Clusters map[string]*InternalCluster
}

func (m *ClustersManager) Add(clusterId string, c *InternalCluster) error {
	m.Lock()
	defer m.Unlock()

	//if c.Monitor == nil {
	//	return fmt.Errorf("add: %s, collector should not be nil", clusterId)
	//}

	_, ok := m.Clusters[clusterId]
	if ok {
		return fmt.Errorf("add: internal cluster %s already exist", clusterId)
	}

	m.Clusters[clusterId] = c

	return nil
}

func (m *ClustersManager) Get(clusterId string) (*InternalCluster, error) {
	m.RLock()
	defer m.RUnlock()

	c, ok := m.Clusters[clusterId]
	if !ok {
		return nil, fmt.Errorf("get: internal cluster %s not found", clusterId)
	}

	//if c.Monitor.ClusterState == ccev1.ClusterUnormal && !(clusterId == constants.ManagerCluster) {
	//	return c, fmt.Errorf("internal cluster %v is  not normal, wait for recover", clusterId)
	//}

	return c, nil
}

func (m *ClustersManager) Del(clusterId string) error {
	m.Lock()
	defer m.Unlock()

	internalCluster, ok := m.Clusters[clusterId]
	if !ok {
		return fmt.Errorf("delete: internal cluster %s not found", clusterId)
	}

	// stop goroutines inside internal cluster
	close(internalCluster.StopCh)

	delete(m.Clusters, clusterId)

	return nil
}

// FuzzyCluster be exported for test
type FuzzyCluster struct {
	Name   string
	Config *rest.Config
	Client kubernetes.Client
}

// AddInternalCluster build internal cluster of cluster and add it
// to multi cluster manager
func (m *ClustersManager) AddInternalCluster(cluster *ccev1.Cluster) error {
	_, err := m.Get(cluster.Spec.ClusterId)
	if err == nil {
		// return Immediately if active internal cluster exist
		return nil
	} else {
		// create internal cluster relate with cluster cr
		config, err := kubeconfig.LoadKubeConfigFromBytes(cluster.Spec.KubeConfig)
		if err != nil {
			return fmt.Errorf("load kubeconfig failed: %v", err)
		}

		//managerCluster, err := m.Get(constants.ManagerCluster)
		//if err != nil {
		//	return err
		//}

		// allocate mem address to avoid nil
		if cluster.Status.ClusterInfo == nil {
			cluster.Status.ClusterInfo = &ccev1.ClusterInfo{}
		}
		cluster.Status.ClusterInfo.State = new(ccev1.ClusterState)

		c := new(InternalCluster)
		c.StopCh = make(chan struct{})
		c.Config = config
		//c.Monitor = monitor.NewClusterMonitor(cluster.Spec.ClusterId, 0, 0, managerCluster.Client.Direct(), c.StopCh)
		c.Client, err = kubernetes.NewClientFor(config, c.StopCh)
		if err != nil {
			return err
		}

		c.VirtClient, err = kubecli.GetKubevirtClientFromRESTConfig(config)
		if err != nil {
			return err
		}

		err = m.Add(cluster.Spec.ClusterId, c)
		if err != nil {
			return fmt.Errorf("add internal cluster failed: %v", err)
		}
	}

	return nil
}

// 收集指定集群数据
/*func (m *ClustersManager) MonitorFor(ctx context.Context, clusterId string) error {
	c, err := m.Get(clusterId)
	if err != nil {
		return err
	}

	c.Monitor.Once.Do(func() {
		log.Info("Start collector for cluster %v", c.Monitor.ClusterId)

		ctx = exit.SetupCtxWithStop(ctx, c.Monitor.StopCh)

		time.AfterFunc(time.Duration(c.Monitor.InitialDelaySeconds)*time.Second, func() {
			go c.Monitor.Collect(ctx)
		})
	})

	return nil
}*/
