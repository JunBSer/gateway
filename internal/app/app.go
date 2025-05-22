package app

import (
	"context"
	"errors"
	"github.com/JunBSer/gateway/internal/config"
	"github.com/JunBSer/gateway/internal/transport/grpc_gateway"
	"github.com/JunBSer/gateway/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

func MustRun(cfg *config.Config) {
	mainLogger := logger.New("Gateway", cfg.Logger.LogLvl)

	endpoints := grpc_gateway.NewEndpointConfig()

	endpoints.PublicPaths = map[string]struct{}{
		"/v1/auth/login":     {},
		"/v1/auth/register":  {},
		"/v1/auth/refresh":   {},
		"/docs/swagger.json": {},
	}

	endpoints.AdminPrefixes = []string{
		"/v1/admin",
	}

	gw := grpc_gateway.NewGateway(&cfg.GW, mainLogger, endpoints)

	graceCh := make(chan os.Signal, 2)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := gw.Start()

		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				mainLogger.Info(context.Background(), "gRpc server stopped", zap.Error(err))
			} else {
				mainLogger.Error(context.Background(), "Error to start gRPC server", zap.Error(err))
			}
		}
	}()

	sig := <-graceCh
	mainLogger.Info(context.Background(), "Shutting down...", zap.String("signal", sig.String()))

	gw.Stop(context.Background())

}
