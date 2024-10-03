package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/yorubanashi/languages/internal/server"
)

const configPath = "config/config.yaml"

func main() {
	logger := log.Default()
	cfg, err := server.LoadConfig(configPath)
	if err != nil {
		logger.Fatal(err)
	}

	srv := server.New(cfg, logger)
	srv.Register()
	srv.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	select {
	case sig := <-c:
		logger.Println(fmt.Sprintf("Received %s, shutting down server", sig.String()))
		srv.Stop()
	}
}
