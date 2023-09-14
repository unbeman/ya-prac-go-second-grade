package main

import (
	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-go-second-grade/internal/client"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/cli"
	"github.com/unbeman/ya-prac-go-second-grade/internal/client/config"
)

func main() {
	cfg, err := config.GetClientConfig()
	if err != nil {
		log.Fatal(err)
	}
	app, err := client.GetClientApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
	executor := cli.GetExecutor(app)

	p := prompt.New(
		executor.Execute,
		cli.Complete,
		prompt.OptionInputTextColor(prompt.Green),
		prompt.OptionTitle("Password manager"),
		prompt.OptionPrefix(">>>"),
	)

	p.Run()
}
