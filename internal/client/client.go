package client

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/unbeman/ya-prac-go-second-grade/internal/client/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/service"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/storage"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/utils"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/utils/interceptor"
)

var protectedWithAuthMethods = map[string]bool{
	"/pass_keeper.OtpService/OTPGenerate": true,
	"/pass_keeper.OtpService/OTPVerify":   true,
	"/pass_keeper.OtpService/OTPValidate": true,
	"/pass_keeper.OtpService/OTPDisable":  true,
	"/pass_keeper.SyncService/Save":       true,
	"/pass_keeper.SyncService/Load":       true,
}

// ClientApp contains services for cli.
type ClientApp struct {
	Auth *service.AuthService
	F2a  *service.OTPService
	Sync *service.SyncService
}

func GetClientApp(cfg config.AppConfig) (*ClientApp, error) {
	vault, err := storage.GetStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not setup storage: %w", err)
	}

	mem, err := storage.NewMemStore()
	if err != nil {
		return nil, fmt.Errorf("could not setup in memory storage: %w", err)
	}

	tlsCredential, err := utils.LoadTLSCredentials(cfg.Certs)
	if err != nil {
		return nil, fmt.Errorf("could not setup certs: %w", err)
	}

	authConn, err := grpc.Dial(cfg.Address,
		grpc.WithTransportCredentials(tlsCredential),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create gRPC connection: %w", err)
	}

	auth := service.NewAuthService(authConn)
	authInterceptor := interceptor.NewAuthInterceptor(auth, protectedWithAuthMethods)

	protectedConn, err := grpc.Dial(
		cfg.Address,
		grpc.WithTransportCredentials(tlsCredential),
		grpc.WithUnaryInterceptor(authInterceptor.Unary()),
		grpc.WithStreamInterceptor(authInterceptor.Stream()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create gRPC connection: %w", err)
	}

	f2a := service.NewOTPService(protectedConn, auth)
	sync := service.NewSyncService(protectedConn, vault, mem, auth)

	return &ClientApp{Auth: auth, F2a: f2a, Sync: sync}, nil
}

func (a ClientApp) Run() {
}

func (a ClientApp) Stop() {
	a.Sync.StopSync()
	log.Info("App stopped")
}
