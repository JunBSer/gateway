package grpc_gateway

import (
	"context"
	"fmt"
	authpb "github.com/JunBSer/services_proto/gen/go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceRegistrar func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

func (s *GatewayServer) registerServices(ctx context.Context, mux *runtime.ServeMux) error {
	const op = "gateway.register.services"

	authAddr := fmt.Sprintf("%s:%s", s.config.AuthHost, s.config.AuthPort)
	err := RegisterAuthService(context.Background(), mux, authAddr)

	if err != nil {
		s.logger.Error(
			context.Background(),
			"Error registering auth service",
			zap.Error(err),
			zap.String("caller", op),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	//in future add registration of booking

	return nil
}

func RegisterAuthService(ctx context.Context, mux *runtime.ServeMux, addr string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return fmt.Errorf("failed to connect to auth service: %w", err)
	}
	return authpb.RegisterAuthHandler(ctx, mux, conn)
}

//func RegisterBookingService(ctx context.Context, mux *runtime.ServeMux, addr string) error {
//	conn, err := grpc.NewClient(addr)
//
//	if err != nil {
//		return fmt.Errorf("failed to connect to booking service: %w", err)
//	}
//	return bookingpb.RegisterBookingHandler(ctx, mux, conn)
//}
