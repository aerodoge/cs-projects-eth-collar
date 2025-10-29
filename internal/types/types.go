package types

type Config struct {
	Deribit    DeribitConfig    `yaml:"deribit" mapstructure:"deribit"`
	Monitor    MonitorConfig    `yaml:"monitor" mapstructure:"monitor"`
	Prometheus PrometheusConfig `yaml:"prometheus" mapstructure:"prometheus"`
	Log        LogConfig        `yaml:"log" mapstructure:"log"`
}

type DeribitConfig struct {
	APIKey    string `yaml:"api_key" mapstructure:"api_key"`
	APISecret string `yaml:"api_secret" mapstructure:"api_secret"`
	TestNet   bool   `yaml:"test_net" mapstructure:"test_net"`
}

type MonitorConfig struct {
	Interval int    `yaml:"interval_seconds" mapstructure:"interval_seconds"`
	Account  string `yaml:"account" mapstructure:"account"`
}

// PrometheusConfig Prometheus 指标服务配置
type PrometheusConfig struct {
	Enabled     bool              `yaml:"enabled" mapstructure:"enabled"`           // 是否启用 Prometheus 指标服务
	PushGateway PushGatewayConfig `yaml:"push_gateway" mapstructure:"push_gateway"` // PushGateway 配置
}

// PushGatewayConfig PushGateway 配置
type PushGatewayConfig struct {
	URL      string            `yaml:"url" mapstructure:"url"`           // PushGateway 地址
	JobName  string            `yaml:"job_name" mapstructure:"job_name"` // 任务名称
	Instance string            `yaml:"instance" mapstructure:"instance"` // 实例标识
	Labels   map[string]string `yaml:"labels" mapstructure:"labels"`     // 额外的标签
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
