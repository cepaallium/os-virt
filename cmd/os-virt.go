package main

import (
	"os"
	"os-virt/cmd/options"
	"os-virt/cmd/virt"
	"os-virt/pkg/clients"
	"os-virt/pkg/utils/log"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	virtConfig := &options.VirtConfig{}
	err := virtConfig.LoadConfig()
	if err != nil {
		log.Error("init ksc config failed: %v", err)
		os.Exit(0)
	}
	loggerConfig := virtConfig.MakeLoggerCfg()
	// init  logger first
	log.InitLoggerWithOpts(loggerConfig)
	// initialize client set
	clients.InitCceManagerClientSetWithOpts()
	manager, err := virt.NewCtrlMgrWithOpts(virtConfig)
	if err != nil {
		log.Error("create controller manager failed: %v", err)
		os.Exit(0)
	}
	err = manager.Initialize(virtConfig)
	if err != nil {
		log.Error("ksc initialized failed: %v", err)
		os.Exit(0)
	}

	if err = manager.CtrlMgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error("unable to run manager", err)
		os.Exit(1)
	}
}
