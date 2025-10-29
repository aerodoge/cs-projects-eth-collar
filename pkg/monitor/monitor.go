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
	// 首先进行 Deribit API 认证
	s.logger.Info("Authenticating with Deribit API")
	if err := s.deribitClient.Authenticate(); err != nil {
		return fmt.Errorf("failed to authenticate with Deribit: %w", err)
	}
	s.logger.Info("Successfully authenticated with Deribit API")

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

	ethEquityUSD := ethSummary.Equity * ethPriceUSD // 计算 ETH 权益的美元价值

	// 计算维持保证金的美元价值
	maintenanceMarginUSD := ethSummary.MaintenanceMargin * ethPriceUSD

	// 计算维持保证金比率（用于日志显示）
	mmRatio := 0.0
	if ethSummary.Equity != 0 {
		mmRatio = ethSummary.MaintenanceMargin / ethSummary.Equity
	}

	// 计算需要补充的ETH数量
	requiredETHAmount := s.calculateRequiredETH(mmRatio, maintenanceMarginUSD, ethSummary.Equity, ethEquityUSD, ethPriceUSD)

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
		zap.Float64("required_eth_amount", requiredETHAmount),
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
		requiredETHAmount,            // 需要补充的ETH数量
		timestamp,                    // 时间戳
	)

	return nil
}

// calculateRequiredETH 计算需要补充的ETH数量
func (s *Service) calculateRequiredETH(mmRatio, maintenanceMarginUSD, ethEquity, ethEquityUSD, ethPriceUSD float64) float64 {
	// 算法1: MM > 50%报警，推送补ETH至MM=30%需要的ETH数量
	if mmRatio > 0.5 {
		// 目标维持保证金比率 = 0.3
		// 0.3 = Current_Total_MM_USD / (Current_Total_Equity_USD + 新增的ETH价值)
		// 新增的ETH价值 = Current_Total_MM_USD / 0.3 - Current_Total_Equity_USD
		targetMMRatio := 0.3
		requiredETHValueUSD := maintenanceMarginUSD/targetMMRatio - ethEquityUSD

		if requiredETHValueUSD > 0 {
			requiredETHAmount := requiredETHValueUSD / ethPriceUSD
			s.logger.Warn("MM ratio alert triggered",
				zap.Float64("current_mm_ratio", mmRatio),
				zap.Float64("target_mm_ratio", targetMMRatio),
				zap.Float64("required_eth_amount", requiredETHAmount),
				zap.Float64("required_eth_value_usd", requiredETHValueUSD),
			)
			return requiredETHAmount
		}
	}

	// 算法2: ETH equity * ETH spot < -0.7m USD报警，补ETH至 ETH equity = 200
	// 当ETH equity为负数时，乘以价格得到负的美元价值，表示亏损
	if ethEquityUSD < -700000 { // -0.7M USD
		requiredETHAmount := 200 - ethEquity
		s.logger.Warn("ETH equity loss alert triggered",
			zap.Float64("current_eth_equity", ethEquity),
			zap.Float64("current_eth_equity_usd", ethEquityUSD),
			zap.Float64("target_eth_equity", 200),
			zap.Float64("required_eth_amount", requiredETHAmount),
			zap.Float64("loss_threshold_usd", -700000),
		)
		return requiredETHAmount
	}

	// 没有触发告警条件，返回0
	return 0
}
