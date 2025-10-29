package metrics

import (
	"cs-projects-eth-collar/internal/types"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"go.uber.org/zap"
)

// Metrics 结构体包含所有 Prometheus 指标和推送配置
type Metrics struct {
	// Prometheus 指标
	MaintenanceMarginRatio *prometheus.GaugeVec // 维持保证金比率指标
	ETHEquity              *prometheus.GaugeVec // ETH权益数量指标
	ETHEquityUSD           *prometheus.GaugeVec // ETH权益美元价值指标
	TotalEquity            *prometheus.GaugeVec // 总权益指标
	MaintenanceMargin      *prometheus.GaugeVec // 维持保证金指标
	MarginBalance          *prometheus.GaugeVec // 保证金余额指标
	ETHPriceUSD            *prometheus.GaugeVec // ETH价格指标
	CollectionTimestamp    *prometheus.GaugeVec // 指标收集时间戳
	RequiredETHAmount      *prometheus.GaugeVec // 需要补充的ETH数量

	// 配置和推送相关
	config   types.PrometheusConfig // Prometheus 配置
	registry *prometheus.Registry   // 指标注册器
	logger   *zap.Logger            // 日志记录器
}

// NewMetrics 创建新的 Metrics 实例，初始化所有 Prometheus 指标
func NewMetrics(config types.PrometheusConfig, logger *zap.Logger) *Metrics {
	// 创建自定义注册器，用于 push 模式
	registry := prometheus.NewRegistry()

	// 创建指标实例
	m := &Metrics{
		config:   config,
		registry: registry,
		logger:   logger,
	}

	// 创建指标并注册到自定义注册器
	m.createMetrics()

	return m
}

// createMetrics 创建指标 (使用自定义注册器推送到 PushGateway)
func (m *Metrics) createMetrics() {
	m.MaintenanceMarginRatio = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deribit_maintenance_margin_ratio",
			Help: "Deribit账户维持保证金比率",
		},
		[]string{"currency", "account"},
	)
	m.ETHEquity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deribit_eth_equity",
			Help: "Deribit账户ETH权益数量",
		},
		[]string{"currency", "account"},
	)
	m.ETHEquityUSD = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deribit_eth_equity_usd",
			Help: "Deribit账户ETH权益美元价值",
		},
		[]string{"currency", "account"},
	)
	m.TotalEquity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deribit_total_equity",
			Help: "Deribit账户总权益",
		},
		[]string{"currency", "account"},
	)
	m.MaintenanceMargin = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deribit_maintenance_margin",
			Help: "Deribit账户维持保证金",
		},
		[]string{"currency", "account"},
	)
	m.MarginBalance = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deribit_margin_balance",
			Help: "Deribit账户保证金余额",
		},
		[]string{"currency", "account"},
	)
	m.ETHPriceUSD = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deribit_eth_price_usd",
			Help: "ETH现货价格(美元)",
		},
		[]string{"currency", "account"},
	)
	m.CollectionTimestamp = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deribit_metrics_collection_timestamp",
			Help: "指标收集的时间戳",
		},
		[]string{"currency", "account"},
	)
	m.RequiredETHAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "deribit_required_eth_amount",
			Help: "需要补充的ETH数量（触发告警时）",
		},
		[]string{"currency", "account"},
	)

	// 注册所有指标到自定义注册器
	m.registry.MustRegister(
		m.MaintenanceMarginRatio,
		m.ETHEquity,
		m.ETHEquityUSD,
		m.TotalEquity,
		m.MaintenanceMargin,
		m.MarginBalance,
		m.ETHPriceUSD,
		m.CollectionTimestamp,
		m.RequiredETHAmount,
	)
}

// UpdateAccountMetrics 更新账户相关的所有 Prometheus 指标
// 参数说明：
//
//	currency: 货币类型 (如 "ETH")
//	account: 账户标识 (如 "default")
//	mmRatio: 维持保证金比率 (0-1 范围)
//	ethEquity: ETH 权益数量
//	ethEquityUSD: ETH 权益美元价值
//	totalEquity: 总权益
//	maintenanceMargin: 维持保证金
//	marginBalance: 保证金余额
//	ethPriceUSD: ETH 现货价格 (美元)
//	requiredETHAmount: 需要补充的ETH数量
//	timestamp: Unix 时间戳
func (m *Metrics) UpdateAccountMetrics(currency, account string, mmRatio, ethEquity, ethEquityUSD, totalEquity, maintenanceMargin, marginBalance, ethPriceUSD, requiredETHAmount float64, timestamp int64) {
	// 创建标签，用于标识不同的货币和账户
	labels := prometheus.Labels{"currency": currency, "account": account} // 指标级别标签

	// 更新各项指标的值
	m.MaintenanceMarginRatio.With(labels).Set(mmRatio)         // 设置维持保证金比率
	m.ETHEquity.With(labels).Set(ethEquity)                    // 设置 ETH 权益数量
	m.ETHEquityUSD.With(labels).Set(ethEquityUSD)              // 设置 ETH 权益美元价值
	m.TotalEquity.With(labels).Set(totalEquity)                // 设置总权益
	m.MaintenanceMargin.With(labels).Set(maintenanceMargin)    // 设置维持保证金
	m.MarginBalance.With(labels).Set(marginBalance)            // 设置保证金余额
	m.ETHPriceUSD.With(labels).Set(ethPriceUSD)                // 设置 ETH 现货价格
	m.RequiredETHAmount.With(labels).Set(requiredETHAmount)    // 设置需要补充的ETH数量
	m.CollectionTimestamp.With(labels).Set(float64(timestamp)) // 设置指标收集时间戳

	// 自动推送指标到 PushGateway
	if err := m.PushMetrics(); err != nil {
		m.logger.Error("Failed to push metrics to PushGateway", zap.Error(err))
	}
}

// PushMetrics 将指标推送到 PushGateway
func (m *Metrics) PushMetrics() error {

	// 创建 Pusher
	pusher := push.New(m.config.PushGateway.URL, m.config.PushGateway.JobName).
		Gatherer(m.registry).
		Grouping("instance", m.config.PushGateway.Instance) // PushGateway级别标签

	// 添加额外的标签
	for key, value := range m.config.PushGateway.Labels {
		pusher = pusher.Grouping(key, value)
	}

	// 推送指标到 PushGateway
	if err := pusher.Push(); err != nil {
		return fmt.Errorf("failed to push metrics: %w", err)
	}

	m.logger.Debug("Successfully pushed metrics to PushGateway",
		zap.String("url", m.config.PushGateway.URL),
		zap.String("job", m.config.PushGateway.JobName),
		zap.String("instance", m.config.PushGateway.Instance),
	)

	return nil
}
