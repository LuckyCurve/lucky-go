# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

lucky-go 是一个基于 Cobra 框架的 Go CLI 工具，提供四大功能模块：云服务管理、SSH连接、金融数据分析和游戏自动化。项目采用模块化架构，每个功能都是独立的子命令。

## Development Commands

### Build and Install
```bash
# 升级所有依赖
make upgrade

# 安装当前项目到 $GOPATH/bin
make install

# 一键执行：升级依赖 + 安装
make all

# 手动构建（调试用）
go build -v ./...
```

### Testing
```bash
# 运行所有测试（带竞态检测）
go test -v -race ./...

# 运行特定模块测试
go test -v ./cloud/...
go test -v ./config/...
go test -v ./server/ssh/...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Code Quality
```bash
# 运行 golangci-lint
golangci-lint run

# 依赖验证
go mod verify
go mod download
```

## Architecture

### Command Structure
- **main.go**: 应用入口，调用 Execute()
- **root.go**: 根命令定义，注册所有子命令模块
- **模块模式**: 每个功能模块实现 `NewCommand() *cobra.Command` 接口

### Core Modules
1. **config**: 配置管理，处理 `~/.lucky-go/config.yaml` 文件
2. **server/ssh**: SSH连接管理，支持直接连接和HTTP服务器模式
3. **cloud**: 腾讯云Lighthouse实例管理（重启功能）
4. **finance**: 金融数据抓取和市盈率计算（基于国债收益率）
5. **game**: Android游戏自动化（ADB相关）

### Configuration System
- 配置文件位置: `~/.lucky-go/config.yaml`
- 配置结构: `map[string]DestinationInstance`
- 每个目标包含: `ssh`、`region`、`instance-id` 字段
- 配置目录不存在时会自动创建

### Key Patterns
- **依赖注入**: SSH模块使用 `var execCommand = exec.Command` 便于测试
- **并发处理**: Finance模块使用goroutine并行获取数据
- **接口抽象**: Finance模块定义 `HTTPClient` 接口便于mock
- **错误处理**: 统一的中文错误消息

### Testing Strategy
- 每个模块都有对应的测试文件
- 使用接口抽象便于单元测试
- CI在多平台和多Go版本下运行测试
- 竞态条件检测是必需的

## Dependencies
- **Cobra**: CLI框架
- **TencentCloud SDK**: 云服务API
- **goquery**: HTML解析（金融数据抓取）
- **yaml.v3**: 配置文件处理

## Environment Notes
- Go 1.24+ required
- 支持多平台编译（Windows/Linux/macOS）
- 金融功能依赖外部网站可用性
- 云功能需要腾讯云凭证配置