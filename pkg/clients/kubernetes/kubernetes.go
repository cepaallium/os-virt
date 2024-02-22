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

package kubernetes

import (
	"context"
	"fmt"
	v1 "os-virt/pkg/api/v1"
	"os-virt/pkg/utils/exit"
	"os-virt/pkg/utils/log"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

)

var (
	scheme = initScheme()
)

func initScheme() *runtime.Scheme {
	s := runtime.NewScheme()
	// cache for all k8s and crd resource
	utilruntime.Must(clientgoscheme.AddToScheme(s))
	//utilruntime.Must(hnc.AddToScheme(s))
	utilruntime.Must(apiextensionsv1.AddToScheme(s))
	utilruntime.Must(v1.AddToScheme(s))
	return s
}

// Client retrieves k8s resource with cache or not
type Client interface {
	Cache() cache.Cache
	Direct() client.Client
	ClientSet() kubernetes.Interface
}

type InternalClient struct {
	client client.Client
	cache  cache.Cache
	rawClientSet kubernetes.Interface
}

// NewClientFor generate client by config
func NewClientFor(cfg *rest.Config, stopCh chan struct{}) (Client, error) {
	var err error
	c := new(InternalClient)

	c.client, err = client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("new k8s client failed: %v", err)
	}

	c.cache, err = cache.New(cfg, cache.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("new k8s cache failed: %v", err)
	}

	c.rawClientSet, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("new raw k8s clientSet failed: %v", err)
	}

	ctx := exit.SetupCtxWithStop(context.Background(), stopCh)

	go func() {
		err = c.cache.Start(ctx)
		if err != nil {
			// that should not happened
			log.Error("start cache failed: %v", err)
		}
	}()
	c.cache.WaitForCacheSync(ctx)

	return c, nil
}

func (c *InternalClient) Cache() cache.Cache {
	return c.cache
}

func (c *InternalClient) Direct() client.Client {
	return c.client
}

func (c *InternalClient) ClientSet() kubernetes.Interface {
	return c.rawClientSet
}
