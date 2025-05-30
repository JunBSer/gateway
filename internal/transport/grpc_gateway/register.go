package grpc_gateway

import (
	"context"
	"fmt"
	authpb "github.com/JunBSer/services_proto/auth/gen/go"
	bookpb "github.com/JunBSer/services_proto/booking/gen/go"
	hotelpb "github.com/JunBSer/services_proto/hotel/gen/go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceRegistrar func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

func (s *GatewayServer) registerServices(ctx context.Context, mux *runtime.ServeMux) error {
	const op = "gateway.register.services"

	authAddr := fmt.Sprintf("%s:%s", s.Config.AuthHost, s.Config.AuthPort)
	err := RegisterAuthService(context.Background(), mux, authAddr)

	if err != nil {
		s.Logger.Error(
			context.Background(),
			"Error registering auth service",
			zap.Error(err),
			zap.String("caller", op),
		)

		return fmt.Errorf("%s: %w", op, err)
	}

	hotelAddr := fmt.Sprintf("%s:%s", s.Config.HotelHost, s.Config.HotelPort)
	err = RegisterHotelService(context.Background(), mux, hotelAddr)
	if err != nil {
		s.Logger.Error(
			context.Background(),
			"Error registering hotel service",
			zap.Error(err),
			zap.String("caller", op),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	bookingAddr := fmt.Sprintf("%s:%s", s.Config.BookingHost, s.Config.BookingPort)
	err = RegisterBookingService(context.Background(), mux, bookingAddr)
	if err != nil {
		s.Logger.Error(
			context.Background(),
			"Error registering booking service",
			zap.Error(err),
			zap.String("caller", op),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func RegisterAuthService(ctx context.Context, mux *runtime.ServeMux, addr string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return fmt.Errorf("failed to connect to auth service: %w", err)
	}
	return authpb.RegisterAuthHandler(ctx, mux, conn)
}

func RegisterHotelService(ctx context.Context, mux *runtime.ServeMux, addr string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return fmt.Errorf("failed to connect to auth service: %w", err)
	}
	return hotelpb.RegisterHotelServiceHandler(ctx, mux, conn)
}

func RegisterBookingService(ctx context.Context, mux *runtime.ServeMux, addr string) error {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return fmt.Errorf("failed to connect to auth service: %w", err)
	}
	return bookpb.RegisterBookingServiceHandler(ctx, mux, conn)
}
