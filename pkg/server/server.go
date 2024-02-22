package server

import (
	"context"
	"fmt"
	"net/http"
	"os-virt/pkg/utils/log"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type ApiServer struct {
	bindAddress string
	port        int
	tlsCertPath *string
	tlsKeyPath  *string
	client      client.Client
}

func NewApiServer(port int, client client.Client) *ApiServer {
	return &ApiServer{
		port:   port,
		client: client,
	}
}

func (server *ApiServer) Start(ctx context.Context) error {
	router := NewRouter(ctx, server.client)
	listenAddress := fmt.Sprintf("%s:%d", server.bindAddress, server.port)
	httpServer := &http.Server{
		Addr:    listenAddress,
		Handler: router,
	}
	//启动apiserver
	go func() {
		if server.tlsCertPath != nil && server.tlsKeyPath != nil {
			log.Info("start https server,address: %s ", listenAddress)
			if err := httpServer.ListenAndServeTLS(*server.tlsCertPath, *server.tlsKeyPath); err != http.ErrServerClosed {
				log.Error("could not start https server", err)
			}
			return
		}
		log.Info("start http server,address: %s ", listenAddress)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("could not start http server", err)
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("shutting down httpserver error", err)
	}
	log.Info("httpserver stopped")

	return nil
}
