# Deribit Position Monitor Makefile

# 变量定义
APP_NAME = monitor
BUILD_DIR = build
BINARY = $(BUILD_DIR)/$(APP_NAME)
CONFIG_FILE = conf/config.yaml
PID_FILE = $(BUILD_DIR)/$(APP_NAME).pid
LOG_FILE = $(BUILD_DIR)/$(APP_NAME).log

# Go 相关变量
GO = go
GOMOD = $(GO) mod
GOBUILD = $(GO) build
GOCLEAN = $(GO) clean
GOTEST = $(GO) test
GOGET = $(GO) get

# 构建标志
LDFLAGS = -ldflags "-s -w"
BUILD_FLAGS = $(LDFLAGS) -o $(BINARY)

# 默认目标
.PHONY: all
all: clean deps build

# 创建构建目录
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# 下载依赖
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# 构建应用程序
.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) ./cmd/monitor
	@echo "✅ 构建完成: $(BINARY)"

# 构建生产版本（优化）
.PHONY: build-prod
build-prod: $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) ./cmd/monitor
	@echo "✅ 生产版本构建完成: $(BINARY)"

# 运行应用程序（前台）
.PHONY: run
run: build
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "❌ 配置文件 $(CONFIG_FILE) 不存在"; \
		exit 1; \
	fi
	$(BINARY) -config $(CONFIG_FILE)

# 后台运行（守护进程）
.PHONY: daemon
daemon: build
	@if [ ! -f $(CONFIG_FILE) ]; then \
		echo "❌ 配置文件 $(CONFIG_FILE) 不存在"; \
		exit 1; \
	fi
	@if [ -f $(PID_FILE) ]; then \
		echo "❌ 服务已在运行，PID: $$(cat $(PID_FILE))"; \
		exit 1; \
	fi
	nohup $(BINARY) -config $(CONFIG_FILE) > $(LOG_FILE) 2>&1 & echo $$! > $(PID_FILE)
	@echo "✅ 服务已启动，PID: $$(cat $(PID_FILE))"
	@echo "📋 日志文件: $(LOG_FILE)"

# 停止守护进程
.PHONY: stop
stop:
	@if [ ! -f $(PID_FILE) ]; then \
		echo "❌ PID 文件不存在，服务可能未运行"; \
		exit 1; \
	fi
	@PID=$$(cat $(PID_FILE)); \
	if ps -p $$PID > /dev/null 2>&1; then \
		kill $$PID && echo "✅ 服务已停止，PID: $$PID"; \
		rm -f $(PID_FILE); \
	else \
		echo "❌ 进程 $$PID 不存在"; \
		rm -f $(PID_FILE); \
	fi

# 重启服务
.PHONY: restart
restart: stop daemon

# 查看服务状态
.PHONY: status
status:
	@if [ -f $(PID_FILE) ]; then \
		PID=$$(cat $(PID_FILE)); \
		if ps -p $$PID > /dev/null 2>&1; then \
			echo "✅ 服务正在运行，PID: $$PID"; \
		else \
			echo "❌ 服务未运行（PID 文件存在但进程不存在）"; \
		fi \
	else \
		echo "❌ 服务未运行"; \
	fi

# 查看日志
.PHONY: logs
logs:
	@if [ -f $(LOG_FILE) ]; then \
		tail -f $(LOG_FILE); \
	else \
		echo "❌ 日志文件不存在: $(LOG_FILE)"; \
	fi

# 查看最近日志
.PHONY: logs-tail
logs-tail:
	@if [ -f $(LOG_FILE) ]; then \
		tail -n 50 $(LOG_FILE); \
	else \
		echo "❌ 日志文件不存在: $(LOG_FILE)"; \
	fi

# 运行测试
.PHONY: test
test:
	$(GOTEST) -v ./...

# 代码格式化
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# 代码检查
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (echo "请安装 golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

# 清理构建文件
.PHONY: clean
clean:
	$(GOCLEAN)
	@if [ -f $(PID_FILE) ]; then \
		echo "⚠️  发现运行中的服务，正在停止..."; \
		$(MAKE) stop; \
	fi
	rm -rf $(BUILD_DIR)
	@echo "✅ 清理完成"

# 安装到系统（需要 root 权限）
.PHONY: install
install: build-prod
	@echo "安装到 /usr/local/bin/$(APP_NAME)"
	sudo cp $(BINARY) /usr/local/bin/$(APP_NAME)
	sudo chmod +x /usr/local/bin/$(APP_NAME)
	@echo "✅ 安装完成"

# 从系统卸载
.PHONY: uninstall
uninstall:
	sudo rm -f /usr/local/bin/$(APP_NAME)
	@echo "✅ 卸载完成"

# 创建发布包
.PHONY: release
release: clean build-prod
	@VERSION=$$(date +%Y%m%d-%H%M%S); \
	RELEASE_DIR="$(BUILD_DIR)/release-$$VERSION"; \
	mkdir -p $$RELEASE_DIR; \
	cp $(BINARY) $$RELEASE_DIR/; \
	cp -r conf $$RELEASE_DIR/; \
	cp README.md $$RELEASE_DIR/; \
	tar -czf $(BUILD_DIR)/$(APP_NAME)-$$VERSION.tar.gz -C $(BUILD_DIR) release-$$VERSION; \
	echo "✅ 发布包已创建: $(BUILD_DIR)/$(APP_NAME)-$$VERSION.tar.gz"

# 开发模式（监听文件变化自动重启）
.PHONY: dev
dev:
	@which air > /dev/null || (echo "请安装 air: go install github.com/cosmtrek/air@latest" && exit 1)
	air

# 显示帮助信息
.PHONY: help
help:
	@echo "Deribit Position Monitor Makefile"
	@echo ""
	@echo "可用命令:"
	@echo "  build       - 构建应用程序"
	@echo "  build-prod  - 构建生产版本（Linux AMD64）"
	@echo "  run         - 前台运行应用程序"
	@echo "  daemon      - 后台运行（守护进程）"
	@echo "  stop        - 停止守护进程"
	@echo "  restart     - 重启服务"
	@echo "  status      - 查看服务状态"
	@echo "  logs        - 实时查看日志"
	@echo "  logs-tail   - 查看最近50行日志"
	@echo "  test        - 运行测试"
	@echo "  fmt         - 格式化代码"
	@echo "  lint        - 代码检查"
	@echo "  clean       - 清理构建文件"
	@echo "  deps        - 下载依赖"
	@echo "  install     - 安装到系统"
	@echo "  uninstall   - 从系统卸载"
	@echo "  release     - 创建发布包"
	@echo "  dev         - 开发模式（需要 air）"
	@echo "  help        - 显示此帮助信息"