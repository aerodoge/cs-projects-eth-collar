package monitor

import (
	"cs-projects-eth-collar/internal/types"
	"cs-projects-eth-collar/pkg/alert"
	"cs-projects-eth-collar/pkg/deribit"
	"fmt"
	"go.uber.org/zap"
	"math"
	"time"
)

type Service struct {
	config        types.MonitorConfig
	deribitClient *deribit.Client
	alertService  *alert.Service
	logger        *zap.Logger
}

func NewService(config types.MonitorConfig, deribitClient *deribit.Client, alertService *alert.Service, logger *zap.Logger) *Service {
	return &Service{
		config:        config,
		deribitClient: deribitClient,
		alertService:  alertService,
		logger:        logger,
	}
}

func (s *Service) Start() error {
	s.logger.Info("Starting position monitor",
		zap.Int("interval_seconds", s.config.Interval),
		zap.Float64("mm_threshold", s.config.MMThreshold),
		zap.Float64("eth_equity_threshold", s.config.ETHEquityThreshold),
	)

	ticker := time.NewTicker(time.Duration(s.config.Interval) * time.Second)
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

	s.logger.Info("Account status check",
		zap.String("currency", "ETH"),
		zap.Float64("equity", ethSummary.Equity),
		zap.Float64("margin_balance", ethSummary.MarginBalance),
		zap.Float64("maintenance_margin", ethSummary.MaintenanceMargin),
	)

	if err := s.checkMaintenanceMargin(ethSummary); err != nil {
		s.logger.Error("MM check failed", zap.Error(err))
	}

	if err := s.checkETHEquity(ethSummary); err != nil {
		s.logger.Error("ETH equity check failed", zap.Error(err))
	}

	return nil
}

func (s *Service) checkMaintenanceMargin(summary *types.AccountSummary) error {
	if summary.Equity == 0 {
		return fmt.Errorf("equity is zero, cannot calculate MM ratio")
	}

	mmRatio := summary.MaintenanceMargin / summary.Equity

	s.logger.Info("MM check",
		zap.Float64("mm_ratio", mmRatio),
		zap.Float64("threshold", s.config.MMThreshold),
		zap.Float64("maintenance_margin", summary.MaintenanceMargin),
		zap.Float64("equity", summary.Equity),
	)

	if mmRatio > s.config.MMThreshold {
		ethNeeded := s.calculateETHNeededForMM(summary)

		alert := types.Alert{
			Type: "MM_THRESHOLD_BREACH",
			Message: fmt.Sprintf("Maintenance Margin ratio %.2f%% exceeds threshold %.2f%%. Need to add %.6f ETH to reach target %.2f%%",
				mmRatio*100, s.config.MMThreshold*100, ethNeeded, s.config.MMTarget*100),
			Currency:     "ETH",
			CurrentValue: mmRatio,
			Threshold:    s.config.MMThreshold,
			Timestamp:    time.Now().Unix(),
		}

		return s.alertService.SendAlert(alert)
	}

	return nil
}

func (s *Service) checkETHEquity(summary *types.AccountSummary) error {
	ethPriceUSD := 3000.0
	ethEquityUSD := summary.Equity * ethPriceUSD

	s.logger.Info("ETH equity check",
		zap.Float64("eth_equity", summary.Equity),
		zap.Float64("eth_equity_usd", ethEquityUSD),
		zap.Float64("threshold_usd", s.config.ETHEquityThreshold),
	)

	if ethEquityUSD < s.config.ETHEquityThreshold {
		ethNeeded := s.config.ETHEquityTarget - summary.Equity

		alert := types.Alert{
			Type: "ETH_EQUITY_THRESHOLD_BREACH",
			Message: fmt.Sprintf("ETH equity $%.2f below threshold $%.2f. Need to add %.6f ETH to reach target %.2f ETH",
				ethEquityUSD, s.config.ETHEquityThreshold, ethNeeded, s.config.ETHEquityTarget),
			Currency:     "ETH",
			CurrentValue: ethEquityUSD,
			Threshold:    s.config.ETHEquityThreshold,
			Timestamp:    time.Now().Unix(),
		}

		return s.alertService.SendAlert(alert)
	}

	return nil
}

func (s *Service) calculateETHNeededForMM(summary *types.AccountSummary) float64 {
	ethPriceUSD := 3000.0

	targetEquityUSD := summary.MaintenanceMargin / s.config.MMTarget
	currentEquityUSD := summary.Equity * ethPriceUSD

	additionalEquityNeeded := targetEquityUSD - currentEquityUSD
	ethNeeded := additionalEquityNeeded / ethPriceUSD

	return math.Max(0, ethNeeded)
}
