package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-second-grade/internal/server"
	"github.com/unbeman/ya-prac-go-second-grade/internal/server/config"
)

func main() {

	cfg, err := config.GetServerConfig()
	if err != nil {
		log.Error(err)
		return
	}

	//todo: setup logger

	app, err := server.GetApp(cfg)
	if err != nil {
		log.Error(err)
		return
	}

	// func waits signal to stop program
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(
			exit,
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
		)

		sig := <-exit
		log.Infof("Got signal '%v'", sig)

		app.Stop()
	}()

	app.Run()
}
