### 命令分解
```
nohup $(BINARY) -config $(CONFIG_FILE) > $(LOG_FILE) 2>&1 & echo $$! > $(PID_FILE)

```
- `nohup`
  - 作用: "no hang up" 的缩写
  - 功能: 让程序在终端关闭后仍然继续运行
  - 原理: 忽略 SIGHUP 信号（挂起信号），这样即使 SSH 连接断开或终端关闭，进程也不会被终止

- `$(BINARY) -config $(CONFIG_FILE)`
  - $(BINARY) = build/eth-collar-monitor
  - 启动应用程序并指定配置文件路径

- `> $(LOG_FILE)`
  - 作用: 重定向标准输出（stdout）到日志文件
  - $(LOG_FILE) = build/eth-collar-monitor.log
  - 程序的正常输出会写入这个文件

- `2>&1`
  - 作用: 重定向标准错误（stderr）到标准输出
  - 2 = stderr 文件描述符
  - &1 = 指向 stdout 的引用
  - 结果: 错误信息也会写入同一个日志文件

- `&`
  - 作用: 在后台运行命令
  - 让 shell 不等待程序执行完成，立即返回控制权

- `echo $$! > $(PID_FILE)`
  - $$! = 上一个后台进程的 PID（进程ID）
  - 将 PID 写入文件 build/eth-collar-monitor.pid
  - 用于后续的进程管理（停止、重启等）

### 整体效果
这个命令的作用是：
- 在后台启动监控程序
- 程序不会因为终端关闭而停止
- 所有输出（包括错误）都记录到日志文件
- 将进程ID保存到文件，方便后续管理

## 使用场景

### 启动守护进程
make daemon

### 查看状态（使用保存的PID）
make status

### 停止服务（通过PID文件找到进程并终止）
make stop

### 查看日志输出
make logs

这是一个典型的 Unix/Linux 守护进程启动模式，常用于生产环境中长时间运行的服务



