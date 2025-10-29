# Deribit 仓位监控系统

一个监控 Deribit 交易账户的系统，用于跟踪维持保证金和 ETH 权益水平，并将指标推送到 Prometheus 进行监控和告警。

## 功能特性

- **维持保证金监控**: 实时监控 MM 比率并推送到 Prometheus
- **ETH 权益监控**: 监控 ETH 权益数量和美元价值
- **Prometheus 集成**: 通过 PushGateway 主动推送指标到 Prometheus
- **可配置监控**: 监控间隔和 Prometheus 端点都可配置

## 配置说明

复制 `conf/config.yaml` 并更新您的 API 凭证：

```yaml
deribit:
  api_key: "YOUR_API_KEY"        # 您的 API 密钥
  api_secret: "YOUR_API_SECRET"  # 您的 API 密钥
  base_url: "https://www.deribit.com/api/v2"
  test_net: false                # 设置为 true 使用测试网

monitor:
  interval_seconds: 30           # 监控间隔（秒）
  account: "default"             # 账户标识

prometheus:
  enabled: true                  # 启用 Prometheus 指标推送
  push_gateway:                  # PushGateway 配置
    url: "http://localhost:9091" # PushGateway 地址
    job_name: "deribit-monitor"  # 任务名称
    instance: "default"          # 实例标识
    labels:                      # 额外的标签
      environment: "production"
      service: "deribit-monitor"
```

## 使用方法

### 构建应用程序:
```bash
go build -o monitor .
```

### 使用默认配置运行:
```bash
./monitor
```

### 使用自定义配置运行:
```bash
./monitor -config /path/to/your/config.yaml
```

## Prometheus 指标

系统会推送以下指标到 Prometheus：

### 指标列表
- `deribit_maintenance_margin_ratio{currency="ETH", account="default"}` - 维持保证金比率
- `deribit_eth_equity{currency="ETH", account="default"}` - ETH 权益数量
- `deribit_eth_equity_usd{currency="ETH", account="default"}` - ETH 权益美元价值
- `deribit_total_equity{currency="ETH", account="default"}` - 总权益
- `deribit_maintenance_margin{currency="ETH", account="default"}` - 维持保证金
- `deribit_margin_balance{currency="ETH", account="default"}` - 保证金余额
- `deribit_eth_price_usd{currency="ETH", account="default"}` - ETH 现货价格 (美元)
- `deribit_metrics_collection_timestamp{currency="ETH", account="default"}` - 指标收集时间戳

### 示例 Prometheus 告警规则
```yaml
groups:
  - name: deribit_alerts
    rules:
      - alert: HighMaintenanceMarginRatio
        expr: deribit_maintenance_margin_ratio > 0.5
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "Deribit 维持保证金比率过高"
          description: "账户 {{ $labels.account }} 的维持保证金比率为 {{ $value | humanizePercentage }}，超过 50% 阈值"

      - alert: LowETHEquity
        expr: deribit_eth_equity_usd < -700000
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Deribit ETH 权益过低"
          description: "账户 {{ $labels.account }} 的 ETH 权益为 ${{ $value | humanize }}，低于 -70万美元阈值"
```

## 使用的 API 端点

- `/private/get_account_summary`: 获取账户权益、保证金和余额信息
- `/private/get_positions`: 获取仓位详情（未来使用）

## 依赖库

- **Viper**: 配置管理
- **Zap**: 结构化日志
- **Prometheus Client**: 指标收集和推送

## 项目结构

```
.
├── cmd/monitor/          # 应用程序入口
├── pkg/
│   ├── config/          # 配置管理
│   ├── deribit/         # Deribit API 客户端
│   ├── metrics/         # Prometheus 指标
│   ├── monitor/         # 监控逻辑
│   └── logger/          # 日志设置
├── internal/types/      # 类型定义
└── conf/               # 配置文件目录
    └── config.yaml     # 配置文件
```

## 安全注意事项

- 安全存储 API 凭证
- 生产环境使用环境变量
- 考虑使用有限权限的 API 密钥
- 开发和测试时启用测试网

## 使用示例

### 启动 PushGateway
```bash
docker run -p 9091:9091 prom/pushgateway
```

### 运行监控程序
```bash
make run
```

### 配置 Prometheus
在 Prometheus 配置文件中添加 PushGateway：
```yaml
scrape_configs:
  - job_name: 'pushgateway'
    static_configs:
      - targets: ['localhost:9091']
    scrape_interval: 30s
```

### 查看推送的指标
访问 PushGateway 界面查看推送的指标：
```
http://localhost:9091/metrics
```

## 当前限制

- 仅支持单一货币监控（仅 ETH）
- 需要外部 Prometheus 和 Alertmanager 进行告警
- 如果 ETH 价格 API 失败，会使用备用价格 $3000