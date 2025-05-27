package grpc_gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/JunBSer/gateway/internal/config"
	"github.com/JunBSer/gateway/internal/metadata"
	"github.com/JunBSer/gateway/pkg/logger"
	pb "github.com/JunBSer/services_proto/auth/gen/go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

type ContextKey string

const IsAdminKey ContextKey = "isa"
const UsIDKey ContextKey = "uid"

type GatewayServer struct {
	Config    *config.Gateway
	Logger    logger.Logger
	Server    *http.Server
	Endpoints *metadata.EndpointConfig
}

func NewGateway(cfg *config.Gateway, logger logger.Logger, endpoints *metadata.EndpointConfig) *GatewayServer {
	if endpoints == nil {
		endpoints = metadata.NewEndpointConfig()
	}
	return &GatewayServer{
		Config:    cfg,
		Logger:    logger,
		Endpoints: endpoints,
	}
}

func (s *GatewayServer) Start() error {
	const op = "gateway.Start"
	ctx := context.Background()

	rootMux := http.NewServeMux()

	if s.Endpoints.IsSwaggerEnabled() {
		s.registerSwagger(rootMux)
	}

	gwMux := runtime.NewServeMux(
		runtime.WithErrorHandler(s.errorHandler),
		runtime.WithRoutingErrorHandler(s.routingErrorHandler),
	)

	if err := s.registerServices(ctx, gwMux); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	handler := s.applyMiddlewares(
		gwMux,
		s.loggingMiddleware(),
		s.corsMiddleware(),
		s.authMiddleware(s.Endpoints),
	)

	rootMux.Handle("/", handler)

	s.Server = &http.Server{
		Addr:    s.Config.Host + ":" + s.Config.Port,
		Handler: rootMux,
	}

	go func() {
		s.Logger.Info(ctx, "Starting gateway server",
			zap.String("addr", s.Server.Addr))

		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Error(ctx, "Server failed",
				zap.Error(err),
				zap.String("caller", op))

		}
	}()

	return nil
}

func (s *GatewayServer) respondError(w http.ResponseWriter, r *http.Request, msg string, code int) {
	s.Logger.Error(r.Context(), "Request error",
		zap.String("path", r.URL.Path),
		zap.Int("code", code),
	)
	http.Error(w, msg, code)
}

func (s *GatewayServer) validateToken(ctx context.Context, token string) (context.Context, bool) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", s.Config.AuthHost, s.Config.AuthPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		s.Logger.Error(ctx, "Cannot connect to auth service", zap.Error(err))
		return ctx, false
	}
	defer conn.Close()

	client := pb.NewAuthClient(conn)
	resp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
	if err != nil || resp == nil || !resp.IsValid {
		s.Logger.Error(ctx, "Token validation failed", zap.Error(err))
		return ctx, false
	}

	ctx = context.WithValue(ctx, UsIDKey, resp.UserId)
	ctx = context.WithValue(ctx, IsAdminKey, resp.IsAdmin)
	return ctx, true
}

func (s *GatewayServer) errorHandler(ctx context.Context, mux *runtime.ServeMux,
	marshaller runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	s.Logger.Error(ctx, "Gateway error",
		zap.Error(err),
		zap.String("path", r.URL.Path),
	)
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaller, w, r, err)
}

func (s *GatewayServer) routingErrorHandler(ctx context.Context, mux *runtime.ServeMux,
	marshaller runtime.Marshaler, w http.ResponseWriter, r *http.Request, code int) {

	s.Logger.Error(ctx, "Routing error",
		zap.Int("code", code),
		zap.String("path", r.URL.Path),
	)
	runtime.DefaultRoutingErrorHandler(ctx, mux, marshaller, w, r, code)
}

func (s *GatewayServer) Stop(ctx context.Context) error {
	s.Logger.Info(ctx, "Shutting down gateway server")
	return s.Server.Shutdown(ctx)
}
