package alert

import (
	"bytes"
	"cs-projects-eth-collar/internal/types"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Service struct {
	config types.AlertsConfig
	logger *zap.Logger
}

func NewService(config types.AlertsConfig, logger *zap.Logger) *Service {
	return &Service{
		config: config,
		logger: logger,
	}
}

func (s *Service) SendAlert(alert types.Alert) error {
	if !s.config.Enabled {
		return nil
	}

	for _, method := range s.config.Methods {
		switch method {
		case "log":
			s.logAlert(alert)
		case "webhook":
			if err := s.sendWebhook(alert); err != nil {
				s.logger.Error("Failed to send webhook alert", zap.Error(err))
			}
		case "email":
			s.logger.Info("Email alerts not implemented yet")
		}
	}

	return nil
}

func (s *Service) logAlert(alert types.Alert) {
	s.logger.Warn("ALERT",
		zap.String("type", alert.Type),
		zap.String("message", alert.Message),
		zap.String("currency", alert.Currency),
		zap.Float64("current_value", alert.CurrentValue),
		zap.Float64("threshold", alert.Threshold),
		zap.Int64("timestamp", alert.Timestamp),
	)
}

func (s *Service) sendWebhook(alert types.Alert) error {
	if s.config.Webhook.URL == "" {
		return fmt.Errorf("webhook URL not configured")
	}

	payload, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(s.config.Webhook.URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}
