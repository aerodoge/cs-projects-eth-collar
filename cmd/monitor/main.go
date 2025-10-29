package main

import (
	"cs-projects-eth-collar/pkg/config"
	"cs-projects-eth-collar/pkg/deribit"
	"cs-projects-eth-collar/pkg/logger"
	"cs-projects-eth-collar/pkg/metrics"
	"cs-projects-eth-collar/pkg/monitor"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	configPath := flag.String("config", "conf/config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 打印配置信息以调试
	log.Printf("Loaded config - Monitor interval: %d seconds, Account: %s", cfg.Monitor.Interval, cfg.Monitor.Account)
	log.Printf("Deribit config - TestNet: %t", cfg.Deribit.TestNet)

	zapLogger, err := logger.NewLogger(cfg.Log)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	// 初始化服务组件
	deribitClient := deribit.NewClient(cfg.Deribit)                 // 创建 Deribit API 客户端
	metricsService := metrics.NewMetrics(cfg.Prometheus, zapLogger) // 创建 Prometheus 指标服务
	monitorService := monitor.NewService(cfg.Monitor, deribitClient, metricsService, zapLogger)

	zapLogger.Info("Starting Deribit position monitor")

	// 启动 Push 模式 - 指标会被推送到 PushGateway
	if cfg.Prometheus.Enabled {
		zapLogger.Info("Metrics will be pushed to PushGateway",
			zap.String("pushgateway_url", cfg.Prometheus.PushGateway.URL),
			zap.String("job_name", cfg.Prometheus.PushGateway.JobName),
		)
	}

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
