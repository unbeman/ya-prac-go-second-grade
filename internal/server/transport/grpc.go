// Package transport contains gRPC server which may be run for apps requests.
package transport

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/unbeman/ya-prac-go-second-grade/api/v1"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/services"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/utils"
)

// GRPCServer contains server entity and additional services.
type GRPCServer struct {
	address     string
	server      *grpc.Server
	authService *services.Auth
	otpService  *services.OTP
	syncService *services.Sync
}

// NewGRPCServer setup server with given services and config settings.
// todo: принимать список сервисов
func NewGRPCServer(cfg config.ServerConfig,
	auth *services.Auth,
	otp *services.OTP,
	sync *services.Sync) (*GRPCServer, error) {
	creds, err := utils.LoadTLSCredentials(cfg.TLS)
	if err != nil {
		return nil, err
	}

	noAuthRequiredMethods := map[string]bool{
		pb.AuthService_Login_FullMethodName:    true,
		pb.AuthService_Register_FullMethodName: true,
	}

	otpMethods := map[string]bool{
		pb.OtpService_OTPValidate_FullMethodName: true,
		pb.OtpService_OTPVerify_FullMethodName:   true,
		pb.OtpService_OTPDisable_FullMethodName:  true,
	}

	ai := NewAuthInterceptor(auth, noAuthRequiredMethods, otpMethods)

	server := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(ai.Unary()),
	)

	return &GRPCServer{
		address:     cfg.Address,
		server:      server,
		authService: auth,
		otpService:  otp,
		syncService: sync,
	}, nil
}

func (g *GRPCServer) GetAddress() string {
	return g.address
}

func (g *GRPCServer) Run() error {
	listen, err := net.Listen("tcp", g.address)
	if err != nil {
		return fmt.Errorf("can't bind address: %w", err)
	}

	pb.RegisterAuthServiceServer(g.server, g.authService)
	pb.RegisterOtpServiceServer(g.server, g.otpService)
	pb.RegisterSyncServiceServer(g.server, g.syncService)

	log.Info("starting gRPC server")
	return g.server.Serve(listen)
}

func (g *GRPCServer) Close() error {
	g.server.GracefulStop()
	return nil
}
