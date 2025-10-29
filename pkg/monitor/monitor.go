package monitor

import (
	"cs-projects-eth-collar/internal/types"
	"cs-projects-eth-collar/pkg/deribit"
	"cs-projects-eth-collar/pkg/metrics"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Service struct {
	config        types.MonitorConfig
	deribitClient *deribit.Client
	metrics       *metrics.Metrics
	logger        *zap.Logger
}

func NewService(config types.MonitorConfig, deribitClient *deribit.Client, metrics *metrics.Metrics, logger *zap.Logger) *Service {
	return &Service{
		config:        config,
		deribitClient: deribitClient,
		metrics:       metrics,
		logger:        logger,
	}
}

func (s *Service) Start() error {
	// 确保间隔时间不为 0，设置最小值为 10 秒
	interval := s.config.Interval
	if interval <= 0 {
		s.logger.Warn("Invalid interval, using default 30 seconds", zap.Int("configured_interval", interval))
		interval = 30
	}

	s.logger.Info("Starting position monitor",
		zap.Int("interval_seconds", interval),
		zap.String("account", s.config.Account),
	)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.checkPositions(); err != nil {
				s.logger.Error("Failed to check positions", zap.Error(err))
			}
		}
	}
}

func (s *Service) checkPositions() error {
	ethSummary, err := s.deribitClient.GetAccountSummary("ETH")
	if err != nil {
		return fmt.Errorf("failed to get ETH account summary: %w", err)
	}

	// 计算指标数据
	timestamp := time.Now().Unix() // 获取当前时间戳

	// 从 Deribit API 获取 ETH 现货价格
	ethPriceUSD, err := s.deribitClient.GetIndexPrice("eth")
	if err != nil {
		s.logger.Error("Failed to get ETH price, using fallback", zap.Error(err))
		ethPriceUSD = 3000.0 // 备用价格
	}

	mmRatio := 0.0 // 维持保证金比率
	if ethSummary.Equity != 0 {
		// 计算维持保证金比率 = 维持保证金 / 权益
		mmRatio = ethSummary.MaintenanceMargin / ethSummary.Equity
	}
	ethEquityUSD := ethSummary.Equity * ethPriceUSD // 计算 ETH 权益的美元价值

	// 记录账户状态信息
	s.logger.Info("Account status check",
		zap.String("currency", "ETH"),
		zap.String("account", s.config.Account),
		zap.Float64("eth_price_usd", ethPriceUSD),
		zap.Float64("equity", ethSummary.Equity),
		zap.Float64("equity_usd", ethEquityUSD),
		zap.Float64("margin_balance", ethSummary.MarginBalance),
		zap.Float64("maintenance_margin", ethSummary.MaintenanceMargin),
		zap.Float64("mm_ratio", mmRatio),
	)

	// 更新 Prometheus 指标
	// 将账户数据推送到 Prometheus，供监控和告警使用
	s.metrics.UpdateAccountMetrics(
		"ETH",                        // 货币类型
		s.config.Account,             // 账户标识
		mmRatio,                      // 维持保证金比率
		ethSummary.Equity,            // ETH 权益数量
		ethEquityUSD,                 // ETH 权益美元价值
		ethSummary.Equity,            // 总权益 (这里与 ETH 权益相同)
		ethSummary.MaintenanceMargin, // 维持保证金
		ethSummary.MarginBalance,     // 保证金余额
		ethPriceUSD,                  // ETH 现货价格
		timestamp,                    // 时间戳
	)

	return nil
}
