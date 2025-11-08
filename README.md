# lucky-go

一个功能强大的 Go 语言命令行工具集合，提供云服务管理、SSH连接、金融数据分析和游戏自动化等功能。

## 功能特性

- ☁️ **云服务管理** - 腾讯云轻量应用服务器实例的重启操作
- 🔐 **SSH连接管理** - 快速连接到配置好的远程服务器
- 📈 **金融数据分析** - 基于国债和AAA公司债券收益率计算市盈率
- 🎮 **游戏自动化** - Android设备的自动点击操作（游戏挂机）

## 安装

### 前置要求

- Go 1.24.0 或更高版本
- 腾讯云账户和访问密钥（用于云服务功能）
- ADB (Android Debug Bridge) 工具（用于游戏自动化功能）

### 从源码安装

```bash
git clone https://github.com/your-username/lucky-go.git
cd lucky-go
make install
```

### 使用 Makefile

```bash
# 升级所有依赖
make upgrade

# 安装当前项目
make install

# 一键执行：升级依赖 + 安装
make all
```

## 配置

配置文件位于：`~/.lucky-go/config.yaml`

示例配置：

```yaml
destinations:
  server1:
    ssh: "user@192.168.1.100"
    region: "ap-beijing"
    instanceId: "lhins-xxxxx"
  server2:
    ssh: "user@192.168.1.101"
    region: "ap-shanghai"
    instanceId: "lhins-yyyyy"
```

## 使用方法

### SSH 服务器管理

```bash
# 连接到指定服务器
lucky-go ssh server1

# 启动 HTTP 服务器
lucky-go ssh serve
```

### 云服务管理

```bash
# 重启指定目标的腾讯云轻量应用服务器
lucky-go cloud reboot server1
```

### 金融数据分析

```bash
# 计算基于债券收益率的市盈率
lucky-go finance pe
```

### 游戏自动化

```bash
# 启动游戏自动点击功能
lucky-go game
```

## 项目结构

```
lucky-go/
├── main.go                    # 应用程序入口点
├── root.go                    # 根命令定义和初始化
├── go.mod                     # Go模块配置
├── Makefile                   # 构建脚本
├── cloud/                     # 云服务模块
│   ├── cli.go                 # 云服务CLI命令
│   └── reboot.go              # 腾讯云实例重启实现
├── config/                    # 配置管理模块
│   └── config.go              # 配置文件处理
├── finance/                   # 金融分析模块
│   └── cli.go                 # 金融计算CLI命令
├── game/                      # 游戏自动化模块
│   └── cli.go                 # 游戏自动化CLI命令
└── server/                    # 服务器管理模块
    └── ssh/                   # SSH连接子模块
        └── cli.go             # SSH连接CLI命令
```

## 技术栈

- **语言**: Go 1.24
- **CLI框架**: [Cobra](https://github.com/spf13/cobra)
- **云服务**: 腾讯云 Lighthouse SDK
- **Web抓取**: [goquery](https://github.com/PuerkitoBio/goquery)
- **配置管理**: YAML
- **自动化**: ADB (Android Debug Bridge)

## 开发

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### 代码质量检查

```bash
# 运行 golangci-lint
golangci-lint run

# 检测竞态条件
go test -race ./...
```

### CI/CD

项目使用 GitHub Actions 进行持续集成和持续部署，支持：

- 多平台测试（Ubuntu、Windows、macOS）
- 代码质量检查
- 竞态条件检测
- 自动化测试流程

## 依赖

主要依赖项：

- `github.com/spf13/cobra` - CLI框架
- `github.com/tencentcloud/tencentcloud-sdk-go` - 腾讯云SDK
- `github.com/PuerkitoBio/goquery` - HTML解析
- `gopkg.in/yaml.v3` - YAML配置处理

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

本项目采用 MIT 许可证。

## 更新日志

### v1.0.0
- 初始版本发布
- 实现SSH连接管理
- 添加腾讯云实例重启功能
- 集成金融数据分析
- 支持游戏自动化功能