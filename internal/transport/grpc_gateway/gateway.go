package grpc_gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/JunBSer/gateway/internal/config"
	"github.com/JunBSer/gateway/pkg/logger"
	pb "github.com/JunBSer/services_proto/gen/go"
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
	config    *config.Gateway
	logger    logger.Logger
	server    *http.Server
	endpoints *EndpointConfig
}

func NewGateway(cfg *config.Gateway, logger logger.Logger, endpoints *EndpointConfig) *GatewayServer {
	if endpoints == nil {
		endpoints = NewEndpointConfig()
	}
	return &GatewayServer{
		config:    cfg,
		logger:    logger,
		endpoints: endpoints,
	}
}

func (s *GatewayServer) Start() error {
	const op = "gateway.Start"
	ctx := context.Background()

	rootMux := http.NewServeMux()

	if s.endpoints.SwaggerEnabled {
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
		s.authMiddleware(),
		s.adminCheckMiddleware(),
	)

	rootMux.Handle("/", handler)

	s.server = &http.Server{
		Addr:    s.config.Host + ":" + s.config.Port,
		Handler: rootMux,
	}

	go func() {
		s.logger.Info(ctx, "Starting gateway server",
			zap.String("addr", s.server.Addr))

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error(ctx, "Server failed",
				zap.Error(err),
				zap.String("caller", op))

		}
	}()

	return nil
}

func (s *GatewayServer) respondError(w http.ResponseWriter, r *http.Request, msg string, code int) {
	s.logger.Error(r.Context(), "Request error",
		zap.String("path", r.URL.Path),
		zap.Int("code", code),
	)
	http.Error(w, msg, code)
}

func (s *GatewayServer) validateToken(ctx context.Context, token string) (context.Context, bool) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", s.config.AuthHost, s.config.AuthPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		s.logger.Error(ctx, "Connection failed", zap.Error(err))
		return ctx, false
	}
	defer conn.Close()

	client := pb.NewAuthClient(conn)

	ctx = context.Background()

	resp, err := client.ValidateToken(ctx, &pb.ValidateTokenRequest{Token: token})
	if err != nil || resp == nil {
		s.logger.Error(ctx, "Validation failed",
			zap.Error(err),
			zap.Bool("nil_response", resp == nil))
		return ctx, false
	}

	if !resp.GetIsValid() || resp.GetUserId().String() == "" {
		s.logger.Error(ctx, "Invalid token data",
			zap.Bool("is_valid", resp.GetIsValid()),
			zap.String("user_id", resp.GetUserId().String()))
		return ctx, false
	}

	return context.WithValue(
		context.WithValue(ctx, UsIDKey, resp.GetUserId()),
		IsAdminKey, resp.GetIsAdmin(),
	), true
}
func (s *GatewayServer) isAdmin(ctx context.Context) bool {
	isAdmin, ok := ctx.Value(IsAdminKey).(bool)
	return ok && isAdmin
}

func (s *GatewayServer) errorHandler(ctx context.Context, mux *runtime.ServeMux,
	marshaller runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	s.logger.Error(ctx, "Gateway error",
		zap.Error(err),
		zap.String("path", r.URL.Path),
	)
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaller, w, r, err)
}

func (s *GatewayServer) routingErrorHandler(ctx context.Context, mux *runtime.ServeMux,
	marshaller runtime.Marshaler, w http.ResponseWriter, r *http.Request, code int) {

	s.logger.Error(ctx, "Routing error",
		zap.Int("code", code),
		zap.String("path", r.URL.Path),
	)
	runtime.DefaultRoutingErrorHandler(ctx, mux, marshaller, w, r, code)
}

func (s *GatewayServer) Stop(ctx context.Context) error {
	s.logger.Info(ctx, "Shutting down gateway server")
	return s.server.Shutdown(ctx)
}
