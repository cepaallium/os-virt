package virt

import (
	"context"
	"fmt"
	"os-virt/cmd/options"
	"os-virt/pkg/scheme"
	"os-virt/pkg/server"
	"os-virt/pkg/utils/constants"
	"os-virt/pkg/utils/exit"
	"os-virt/pkg/utils/log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

type ControllerManager struct {
	Config  *options.VirtConfig
	CtrlMgr ctrl.Manager
}

func NewCtrlMgrWithOpts(mc *options.VirtConfig) (*ControllerManager, error) {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme.Scheme,
		//Port:           mc.WebhookPort,
		LeaderElection: false,
		Metrics: metricsserver.Options{
			BindAddress: ":8080",
		},
		HealthProbeBindAddress:  ":8081",
		LeaderElectionID:        "os-virt",
		LeaderElectionNamespace: constants.NameSpaceRoot,
	})

	if err != nil {
		log.Fatal("unable to set up controller manager: %v", err)
		return nil, err
	}

	return &ControllerManager{Config: mc, CtrlMgr: mgr}, nil
}

func (m *ControllerManager) Initialize(config *options.VirtConfig) error {
	if err := m.CtrlMgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up health check: %s", err)
	}
	if err := m.CtrlMgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up ready check: %s", err)
	}
	m.CtrlMgr.Add(server.NewApiServer(m.Config.WebserverPort, m.CtrlMgr.GetClient()))
	return nil
}

func (m *ControllerManager) Run(stop <-chan struct{}) {
	err := m.CtrlMgr.Start(exit.SetupCtxWithStop(context.Background(), stop))
	if err != nil {
		log.Fatal("problem run controller manager: %v", err)
	}
}
