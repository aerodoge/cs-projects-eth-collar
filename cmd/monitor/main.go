package main

import (
	"cs-projects-eth-collar/pkg/alert"
	"cs-projects-eth-collar/pkg/config"
	"cs-projects-eth-collar/pkg/deribit"
	"cs-projects-eth-collar/pkg/logger"
	"cs-projects-eth-collar/pkg/monitor"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configPath := flag.String("config", "conf/config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	zapLogger, err := logger.NewLogger(cfg.Log)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	deribitClient := deribit.NewClient(cfg.Deribit)
	alertService := alert.NewService(cfg.Alerts, zapLogger)
	monitorService := monitor.NewService(cfg.Monitor, deribitClient, alertService, zapLogger)

	zapLogger.Info("Starting Deribit position monitor")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := monitorService.Start(); err != nil {
			zapLogger.Fatal("Monitor service failed")
		}
	}()

	<-c
	zapLogger.Info("Shutting down monitor")
}
