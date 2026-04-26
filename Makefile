BINARY  := cex-price-monitoring
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: build build-linux pack clean run dev

# 当前平台（macOS 本地测试用）
build:
	go build $(LDFLAGS) -o $(BINARY) .

# Linux x86 生产包
build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY) .

# 编译 Linux 可执行文件（直接用于部署）
pack: build-linux
	@echo "✅ 编译完成: $(BINARY)  (Linux amd64)"

clean:
	rm -f $(BINARY)

# 本地运行（生产配置）
run: build
	./$(BINARY)

# 本地运行（开发配置）
dev: build
	APP_ENV=dev ./$(BINARY)
