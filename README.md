# Deribit 仓位监控系统

一个监控 Deribit 交易账户的系统，用于跟踪维持保证金和 ETH 权益水平，并在阈值被突破时提供警报。

## 功能特性

- **维持保证金监控**: 当 MM 比率超过 50% 时发出警报，计算达到 30% 目标所需的 ETH 数量
- **ETH 权益监控**: 当 ETH 权益低于 -70 万美元时发出警报，计算达到 200 ETH 目标所需的 ETH 数量
- **多种警报方式**: 日志、Webhook 和邮件（邮件功能尚未实现）
- **可配置阈值**: 所有监控参数都可通过 YAML 配置

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
  mm_threshold: 0.5             # 50% MM 阈值触发警报
  mm_target: 0.3                # 添加 ETH 后的 30% 目标 MM
  eth_equity_threshold: -700000.0  # -70万美元 ETH 权益阈值
  eth_equity_target: 200.0         # 目标 ETH 权益

alerts:
  enabled: true                  # 启用警报
  methods: ["log", "webhook"]    # 警报方式
  webhook:
    url: "https://your-webhook-url.com/alert"  # Webhook 地址
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

## 警报逻辑

### 维持保证金警报
- **触发条件**: 当 MM 比率 > 50%
- **计算公式**:
  ```
  新保证金率 = 0.3
  0.3 = Current_Total_MM_USD / (Current_Total_Equity_USD + 新增的ETH价值)
  新增的ETH数量 = 新增的ETH价值 / ETH价格
  ```

### ETH 权益警报
- **触发条件**: 当 ETH equity * ETH spot price < -70万美元
- **计算公式**:
  ```
  新增ETH的数量 = 200 - ETH equity
  ```

## 使用的 API 端点

- `/private/get_account_summary`: 获取账户权益、保证金和余额信息
- `/private/get_positions`: 获取仓位详情（未来使用）

## 依赖库

- **Viper**: 配置管理
- **Zap**: 结构化日志
- **Gorilla WebSocket**: 用于潜在的实时更新（未来使用）

## 项目结构

```
.
├── cmd/monitor/          # 应用程序入口
├── pkg/
│   ├── config/          # 配置管理
│   ├── deribit/         # Deribit API 客户端
│   ├── alert/           # 警报系统
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

## 当前限制

- ETH 价格当前硬编码为 $3000（应从 API 获取）
- 邮件警报尚未实现
- 仅支持单一货币监控（仅 ETH）
- 无警报历史持久化功能