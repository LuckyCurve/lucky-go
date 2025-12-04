# CODEBUDDY.md

This file provides guidance to CodeBuddy Code when working with code in this repository.

## Project Overview

lucky-go 是一个基于 Cobra 框架的 Go CLI 工具，提供六大功能模块：云服务管理、SSH连接、金融数据分析、Telegram通知推送、汇率查询和游戏自动化。

## Development Commands

```bash
# 构建项目
go build -v ./...

# 安装到 $GOPATH/bin
make install

# 升级依赖并安装
make all

# 运行所有测试（带竞态检测）
go test -v -race ./...

# 运行特定模块测试
go test -v ./finance/...
go test -v ./notify/...
go test -v ./forex/...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 代码质量检查
golangci-lint run
```

## Architecture

### Command Registration Pattern

所有模块在 `root.go` 中注册，每个模块实现 `NewCommand() *cobra.Command` 接口：

```go
// root.go
rootCmd.AddCommand(ssh.NewCommand())
rootCmd.AddCommand(cloud.NewCommand())
rootCmd.AddCommand(finance.NewCommand())
rootCmd.AddCommand(forex.NewCommand())
rootCmd.AddCommand(game.NewCommand())
```

### Module Structure

```
├── config/           # 配置管理 - 处理 ~/.lucky-go/config.yaml
├── cloud/            # 腾讯云Lighthouse实例管理
├── finance/          # FRED API 金融数据获取和PE计算，支持Telegram推送
├── notify/           # Telegram消息推送底层实现
├── forex/            # 汇率查询（Frankfurter API，依赖notify）
├── valuation/        # 标普500 CAPE 估值（Multpl.com 数据）
├── game/             # Android ADB游戏自动化
└── server/ssh/       # SSH连接管理
```

### Module Dependencies

```
forex     ──→ notify ──→ Telegram API
finance   ──→ notify ──→ Telegram API
valuation ──→ finance ──→ FRED API
valuation ──→ notify  ──→ Telegram API
cloud     ──→ config  ──→ ~/.lucky-go/config.yaml
ssh       ──→ config
game      ──→ (独立，仅依赖ADB)
```

### Testability Pattern

各模块使用函数变量注入实现可测试性：

```go
// 生产代码
var execCommand = exec.Command
var rebootInstanceFunc = defaultRebootInstance
var defaultHTTPClient HTTPClient = &http.Client{}

// 测试中替换
execCommand = fakeExecCommand
rebootInstanceFunc = mockRebootFunc
```

### Configuration Structure

配置文件位于 `~/.lucky-go/config.yaml`：

```yaml
dest:
  server1:
    ssh: "user@host"
    region: "ap-beijing"
    instance-id: "lhins-xxxxx"
```

## Environment Variables

| 变量 | 模块 | 用途 |
|------|------|------|
| FRED_API_KEY | finance | FRED API 密钥 |
| TELEGRAM_BOT_TOKEN | notify | Telegram Bot Token |
| TELEGRAM_CHAT_ID | notify | 目标 Chat ID |
| TENCENT_CLOUD_SECRET_ID | cloud | 腾讯云密钥 ID |
| TENCENT_CLOUD_SECRET_KEY | cloud | 腾讯云密钥 |

## Key Implementation Notes

### Telegram Markdown Limitations

Telegram 的 Markdown 模式不支持表格语法 `| col |`。使用代码块 ` ``` ` 配合等宽字体实现表格对齐效果。

### Concurrent Data Fetching

finance 和 notify 模块使用 goroutine + channel 并行获取多个数据源：

```go
treasuryCh := make(chan result, 1)
go func() {
    value, err := finance.Get10YearTreasuryYield()
    treasuryCh <- result{value: value, err: err}
}()
```

### CLI Commands

```
lucky-go
├── cloud reboot [dest]           # 重启腾讯云实例
├── pe                            # 显示PE估值表格
│   └── --push, -p                # 推送结果到Telegram
├── cape                          # 查询标普500 CAPE 估值
│   └── --push, -p                # 推送结果到Telegram
├── forex [from] [to]             # 查询汇率（如 forex USD CNY）
│   └── --amount, -a              # 兑换金额
│   └── --push, -p                # 推送结果到Telegram
├── ssh [dest]                    # SSH连接服务器
│   └── serve --port PORT         # 启动HTTP服务
└── game                          # 启动游戏自动点击
```

## Dependencies

- **Cobra**: CLI框架
- **TencentCloud SDK**: 云服务API
- **tablewriter + fatih/color**: 终端表格渲染和彩色输出
- **yaml.v3**: 配置文件处理

## CI/CD

项目使用 GitHub Actions，支持多平台（Ubuntu/Windows/macOS）和多Go版本（1.21-1.25）测试。
