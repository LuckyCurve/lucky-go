# lucky-go 项目说明

## 项目概述

lucky-go 是一个功能强大的 Go 语言命令行工具集合，提供云服务管理、SSH连接、金融数据分析和游戏自动化等功能。该项目采用模块化架构，基于 Cobra CLI 框架构建，支持多种实用功能，包括：

- ☁️ **云服务管理** - 腾讯云轻量应用服务器实例的重启操作
- 🔐 **SSH连接管理** - 快速连接到配置好的远程服务器
- 📈 **金融数据分析** - 基于国债和AAA公司债券收益率计算市盈率
- 🎮 **游戏自动化** - Android设备的自动点击操作（游戏挂机）

## 技术栈

- **语言**: Go 1.24+
- **CLI框架**: [Cobra](https://github.com/spf13/cobra)
- **云服务**: 腾讯云 Lighthouse SDK
- **Web抓取**: [goquery](https://github.com/PuerkitoBio/goquery)
- **配置管理**: YAML
- **自动化**: ADB (Android Debug Bridge)

## 项目架构

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
├── server/                    # 服务器管理模块
│   └── ssh/                   # SSH连接子模块
│       └── cli.go             # SSH连接CLI命令
└── test files...              # 测试文件
```

## 主要依赖

- `github.com/spf13/cobra` - CLI框架
- `github.com/tencentcloud/tencentcloud-sdk-go` - 腾讯云SDK
- `github.com/PuerkitoBio/goquery` - HTML解析
- `gopkg.in/yaml.v3` - YAML配置处理

## 功能模块分析

### 配置模块 (`config/`)
- 管理用户配置，存储在 `~/.lucky-go/config.yaml`
- 支持多个目标实例的定义，包括SSH连接信息、云区域和实例ID

### SSH模块 (`server/ssh/`)
- 提供与远程服务器的SSH连接功能
- 支持通过配置的目标名称快速连接
- 包含HTTP服务器功能，可指定端口

### 云服务模块 (`cloud/`)
- 提供腾讯云轻量应用服务器的管理功能
- 目前实现实例重启功能
- 通过环境变量读取腾讯云凭证（TENCENT_CLOUD_SECRET_ID, TENCENT_CLOUD_SECRET_KEY）

### 金融分析模块 (`finance/`)
- 从外部源获取10年期国债和AAA公司债券收益率
- 计算各种百分比下的市盈率值
- 使用并行请求提高数据获取效率

### 游戏自动化模块 (`game/`)
- 提供Android设备的自动化点击功能
- 使用ADB工具与设备交互
- 支持多设备选择和自动点击

## 建立和运行

### 前置要求

- Go 1.24.0 或更高版本
- 腾讯云账户和访问密钥（用于云服务功能）
- ADB (Android Debug Bridge) 工具（用于游戏自动化功能）

### 构建命令

```bash
# 从源码安装
make install

# 升级所有依赖
make upgrade

# 一键执行：升级依赖 + 安装
make all
```

### 使用方式

```bash
# SSH 服务器管理
lucky-go ssh server1

# 启动 HTTP 服务器
lucky-go ssh serve

# 云服务管理 - 重启指定目标的腾讯云轻量应用服务器
lucky-go cloud reboot server1

# 金融数据分析 - 计算基于债券收益率的市盈率
lucky-go finance pe

# 游戏自动化 - 启动游戏自动点击功能
lucky-go game
```

## 开发约定

### 代码结构
- 采用模块化设计，每个功能模块独立放置在专门的目录中
- 使用Cobra框架处理CLI命令和子命令
- 配置管理模块统一处理应用配置

### 测试实践
- 项目包含单元测试，主要在 `_test.go` 文件中
- 使用函数变量来支持测试中的模拟功能
- 测试覆盖主要功能但注意Execute函数的测试限制（因调用os.Exit）

### 错误处理
- 使用Go的标准错误处理模式
- 在CLI命令中返回错误以正确处理用户输入

### 配置管理
- 使用YAML格式的配置文件
- 遵循约定的目录结构 `~/.lucky-go/config.yaml`
- 提供配置文件不存在时的自动创建功能

## 部署与发布

项目使用 GitHub Actions 进行持续集成和持续部署，支持：
- 多平台测试（Ubuntu、Windows、macOS）
- 代码质量检查
- 竞态条件检测
- 自动化测试流程

## 扩展性

项目采用模块化设计，易于扩展新功能：
1. 创建新的功能模块目录
2. 实现CLI命令接口
3. 在 `root.go` 中注册新命令
4. 遵循现有代码的错误处理和测试模式

## 维护重点

- 确保外部API调用的错误处理和重试机制
- 维护配置验证逻辑
- 定期更新依赖库以保持安全性和功能完整性
- 持续改进测试覆盖范围，特别是对关键功能的测试

## 开发工作流

- 使用Go Modules管理依赖
- 遵循Go语言最佳实践
- 使用golangci-lint进行代码质量检查
- 编写单元测试并使用覆盖率工具验证测试效果