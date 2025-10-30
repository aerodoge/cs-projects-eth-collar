# Deribit 仓位监控系统

一个监控 Deribit 交易账户的系统，用于跟踪维持保证金和 ETH 权益水平，并将指标推送到 Prometheus 进行监控和告警。

## 功能特性

- **维持保证金监控**: 实时监控 MM 比率并推送到 Prometheus
- **ETH 权益监控**: 监控 ETH 权益数量和美元价值
- **Prometheus 集成**: 通过 PushGateway 主动推送指标到 Prometheus
- **可配置监控**: 监控间隔和 Prometheus 端点都可配置
- **结构化日志**: 基于 Zap 的高性能日志记录
- **守护进程模式**: 支持后台运行和进程管理

## 配置说明

复制 `config.yaml.example` 到 `conf/config.yaml` 并更新您的 API 凭证：

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

log:
  level: "info"                  # 日志级别
  file: "monitor.log"            # 日志文件
```

## 使用方法

### 快速开始

#### 1. 配置文件设置
```bash
# 复制配置文件模板
cp config.yaml.example conf/config.yaml

# 编辑配置文件，填入您的 Deribit API 凭证
vim conf/config.yaml
```

#### 2. 构建和运行

使用 Makefile 构建和运行：

```bash
# 构建应用程序
make build

# 前台运行（开发/调试）
make run

# 后台运行（生产环境）
make daemon

# 查看服务状态
make status

# 查看日志
make logs

# 停止服务
make stop
```

### 手动构建和运行

#### 构建应用程序:
```bash
go build -o build/monitor ./cmd/monitor
```

#### 使用默认配置运行:
```bash
./build/monitor
```

#### 使用自定义配置运行:
```bash
./build/monitor -config /path/to/your/config.yaml
```

### 其他常用命令

```bash
# 安装依赖
make deps

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make lint

# 创建发布包
make release

# 开发模式（热重载，需要安装 air）
make dev

# 查看所有可用命令
make help
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

- **Viper**: 配置管理和解析
- **Zap**: 高性能结构化日志
- **Prometheus Client**: 指标收集和推送
- **Testify**: 单元测试框架

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

## 部署和运维

### Docker 部署 PushGateway
```bash
# 启动 PushGateway
docker run -d -p 9091:9091 --name pushgateway prom/pushgateway

# 查看 PushGateway 状态
docker ps | grep pushgateway
```

### 生产环境部署
```bash
# 构建生产版本
make build-prod

# 安装到系统路径
sudo make install

# 创建系统服务（可选）
sudo systemctl enable deribit-monitor
sudo systemctl start deribit-monitor
```

### 配置 Prometheus
在 Prometheus 配置文件中添加 PushGateway：
```yaml
scrape_configs:
  - job_name: 'pushgateway'
    static_configs:
      - targets: ['localhost:9091']
    scrape_interval: 30s
    honor_labels: true
```

### 监控和调试
```bash
# 查看实时日志
make logs

# 查看最近日志
make logs-tail

# 查看服务状态
make status

# 重启服务
make restart

# 查看推送的指标
curl http://localhost:9091/metrics
```

## 开发和测试

### 运行单元测试
```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./pkg/deribit

# 运行测试并生成覆盖率报告
go test -cover ./...
```

### 开发模式
```bash
# 安装 air（热重载工具）
go install github.com/cosmtrek/air@latest

# 启动开发模式
make dev
```

### 代码质量检查
```bash
# 格式化代码
make fmt

# 静态代码分析
make lint
```

## 故障排除

### 常见问题

1. **API 认证失败**
   - 检查 API 密钥和密钥是否正确
   - 确认是否使用了正确的网络（测试网/主网）

2. **PushGateway 连接失败**
   - 确认 PushGateway 是否正常运行
   - 检查网络连接和防火墙设置

3. **配置文件错误**
   - 验证 YAML 语法是否正确
   - 检查必需字段是否都已填写

### 日志级别
支持的日志级别：`debug`, `info`, `warn`, `error`, `panic`, `fatal`

```yaml
log:
  level: "debug"  # 开发时使用 debug，生产时使用 info
```

## 当前限制

- 仅支持单一货币监控（ETH）
- 需要外部 Prometheus 和 Alertmanager 进行告警
- 如果 ETH 价格 API 失败，会使用备用价格 $3000

## 版本历史

- **v1.0.0**: 基础监控功能，支持维持保证金和权益监控
- **v1.1.0**: 添加结构化日志和守护进程支持
- **v1.2.0**: 增加单元测试和代码质量检查