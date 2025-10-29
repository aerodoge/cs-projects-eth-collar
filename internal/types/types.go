package types

type Config struct {
	Deribit DeribitConfig `yaml:"deribit"`
	Monitor MonitorConfig `yaml:"monitor"`
	Alerts  AlertsConfig  `yaml:"alerts"`
	Log     LogConfig     `yaml:"log"`
}

type DeribitConfig struct {
	APIKey    string `yaml:"api_key"`
	APISecret string `yaml:"api_secret"`
	BaseURL   string `yaml:"base_url"`
	TestNet   bool   `yaml:"test_net"`
}

type MonitorConfig struct {
	Interval           int     `yaml:"interval_seconds"`
	MMThreshold        float64 `yaml:"mm_threshold"`
	MMTarget           float64 `yaml:"mm_target"`
	ETHEquityThreshold float64 `yaml:"eth_equity_threshold"`
	ETHEquityTarget    float64 `yaml:"eth_equity_target"`
}

type AlertsConfig struct {
	Enabled bool          `yaml:"enabled"`
	Methods []string      `yaml:"methods"`
	Webhook WebhookConfig `yaml:"webhook"`
	Email   EmailConfig   `yaml:"email"`
}

type WebhookConfig struct {
	URL string `yaml:"url"`
}

type EmailConfig struct {
	SMTP     string `yaml:"smtp"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type AccountSummary struct {
	Currency          string  `json:"currency"`
	Equity            float64 `json:"equity"`
	MarginBalance     float64 `json:"margin_balance"`
	MaintenanceMargin float64 `json:"maintenance_margin"`
	InitialMargin     float64 `json:"initial_margin"`
	TotalPL           float64 `json:"total_pl"`
	SessionUPL        float64 `json:"session_upl"`
}

type Alert struct {
	Type         string  `json:"type"`
	Message      string  `json:"message"`
	Currency     string  `json:"currency"`
	CurrentValue float64 `json:"current_value"`
	Threshold    float64 `json:"threshold"`
	Timestamp    int64   `json:"timestamp"`
}
