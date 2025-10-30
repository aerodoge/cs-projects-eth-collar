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

type Limits struct {
	MatchingEngine    MatchingEngineLimits `json:"matching_engine"`
	LimitsPerCurrency bool                 `json:"limits_per_currency"`
	NonMatchingEngine LimitConfig          `json:"non_matching_engine"`
}

type MatchingEngineLimits struct {
	BlockRfqMaker        LimitConfig `json:"block_rfq_maker"`
	CancelAll            LimitConfig `json:"cancel_all"`
	GuaranteedMassQuotes LimitConfig `json:"guaranteed_mass_quotes"`
	MaximumMassQuotes    LimitConfig `json:"maximum_mass_quotes"`
	MaximumQuotes        LimitConfig `json:"maximum_quotes"`
	Spot                 LimitConfig `json:"spot"`
	Trading              struct {
		Total LimitConfig `json:"total"`
	} `json:"trading"`
}

type LimitConfig struct {
	Rate  int `json:"rate"`
	Burst int `json:"burst"`
}

type FeeConfig struct {
	Value struct {
		Default struct {
			Type  string  `json:"type"`
			Taker float64 `json:"taker"`
			Maker float64 `json:"maker"`
		} `json:"default"`
		BlockTrade float64 `json:"block_trade"`
	} `json:"value"`
	Kind      string `json:"kind"`
	IndexName string `json:"index_name"`
}

type ChangeMarginModelAPILimit struct {
	Timeframe int64 `json:"timeframe"`
	Rate      int   `json:"rate"`
}

type CurrencySummary struct {
	MarginBalance                float64            `json:"margin_balance"`
	TotalMaintenanceMarginUSD    float64            `json:"total_maintenance_margin_usd,omitempty"`
	OptionsSessionUpl            float64            `json:"options_session_upl"`
	AvailableWithdrawalFunds     float64            `json:"available_withdrawal_funds"`
	ProjectedDeltaTotal          float64            `json:"projected_delta_total"`
	SessionRpl                   float64            `json:"session_rpl"`
	ProjectedInitialMargin       float64            `json:"projected_initial_margin"`
	Limits                       Limits             `json:"limits"`
	OptionsVega                  float64            `json:"options_vega"`
	DepositAddress               string             `json:"deposit_address,omitempty"`
	OptionsSessionRpl            float64            `json:"options_session_rpl"`
	OptionsGammaMap              map[string]float64 `json:"options_gamma_map"`
	AvailableFunds               float64            `json:"available_funds"`
	TotalMarginBalanceUSD        float64            `json:"total_margin_balance_usd,omitempty"`
	TotalPl                      float64            `json:"total_pl"`
	OptionsGamma                 float64            `json:"options_gamma"`
	FuturesSessionUpl            float64            `json:"futures_session_upl"`
	TotalDeltaTotalUSD           float64            `json:"total_delta_total_usd,omitempty"`
	SessionUpl                   float64            `json:"session_upl"`
	OptionsValue                 float64            `json:"options_value"`
	ProjectedMaintenanceMargin   float64            `json:"projected_maintenance_margin"`
	MaintenanceMargin            float64            `json:"maintenance_margin"`
	TotalInitialMarginUSD        float64            `json:"total_initial_margin_usd,omitempty"`
	OptionsVegaMap               map[string]float64 `json:"options_vega_map"`
	OptionsThetaMap              map[string]float64 `json:"options_theta_map"`
	CrossCollateralEnabled       bool               `json:"cross_collateral_enabled"`
	Equity                       float64            `json:"equity"`
	MarginModel                  string             `json:"margin_model"`
	FeeGroup                     string             `json:"fee_group,omitempty"`
	InitialMargin                float64            `json:"initial_margin"`
	FuturesPl                    float64            `json:"futures_pl"`
	Balance                      float64            `json:"balance"`
	AdditionalReserve            float64            `json:"additional_reserve"`
	Currency                     string             `json:"currency"`
	Fees                         []FeeConfig        `json:"fees,omitempty"`
	PortfolioMarginingEnabled    bool               `json:"portfolio_margining_enabled"`
	DeltaTotalMap                map[string]float64 `json:"delta_total_map"`
	OptionsTheta                 float64            `json:"options_theta"`
	TotalEquityUSD               float64            `json:"total_equity_usd,omitempty"`
	SpotReserve                  float64            `json:"spot_reserve"`
	DeltaTotal                   float64            `json:"delta_total"`
	OptionsPl                    float64            `json:"options_pl"`
	OptionsDelta                 float64            `json:"options_delta"`
	FuturesSessionRpl            float64            `json:"futures_session_rpl"`
	FeeBalance                   float64            `json:"fee_balance"`
	LockedBalance                float64            `json:"locked_balance"`
	EstimatedLiquidationRatio    float64            `json:"estimated_liquidation_ratio,omitempty"`
	EstimatedLiquidationRatioMap map[string]float64 `json:"estimated_liquidation_ratio_map,omitempty"`
}

type AccountSummaries struct {
	ID                               int                       `json:"id,omitempty"`
	Type                             string                    `json:"type,omitempty"`
	MandatoryTfa                     bool                      `json:"mandatory_tfa,omitempty"`
	Email                            string                    `json:"email,omitempty"`
	Username                         string                    `json:"username,omitempty"`
	BlockRfqSelfMatchPrevention      bool                      `json:"block_rfq_self_match_prevention,omitempty"`
	CreationTimestamp                int64                     `json:"creation_timestamp,omitempty"`
	SecurityKeysEnabled              bool                      `json:"security_keys_enabled,omitempty"`
	SystemName                       string                    `json:"system_name,omitempty"`
	MmpEnabled                       bool                      `json:"mmp_enabled,omitempty"`
	SelfTradingExtendedToSubaccounts bool                      `json:"self_trading_extended_to_subaccounts,omitempty"`
	ChangeMarginModelAPILimit        ChangeMarginModelAPILimit `json:"change_margin_model_api_limit"`
	Summaries                        []CurrencySummary         `json:"summaries"`
	SelfTradingRejectMode            string                    `json:"self_trading_reject_mode,omitempty"`
	InteruserTransfersEnabled        bool                      `json:"interuser_transfers_enabled,omitempty"`
	ReferrerID                       string                    `json:"referrer_id,omitempty"`
}
