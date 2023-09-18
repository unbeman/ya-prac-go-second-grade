package server

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/database"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/services"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/transport"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/utils"
)

type App struct {
	server *transport.GRPCServer
}

func GetApp(cfg config.ServerConfig) (*App, error) {
	db, err := database.GetDatabase(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("coud not get database: %w", err)
	}

	jwt, err := utils.NewJWTManager(cfg.JWT)
	if err != nil {
		return nil, fmt.Errorf("coud not get jwt manager: %w", err)
	}

	// setup grpc services
	authService := services.NewAuthService(db, jwt)
	otpService := services.NewOTPService(cfg.OTP, db, jwt)
	syncService := services.NewSyncService(db)

	gServer, err := transport.NewGRPCServer(cfg, authService, otpService, syncService)
	if err != nil {
		return nil, fmt.Errorf("coud not get grpc server: %w", err)
	}

	return &App{server: gServer}, nil
}

func (a App) Run() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.server.Run()
		if err != nil {
			log.Info("gRPC server stopped, reason: ", err)
		}
		log.Info("gRPC server stopped successfully")
	}()

	wg.Wait()
}

func (a App) Stop() {
	log.Info("stopping app processes")
	a.server.Close()
}
