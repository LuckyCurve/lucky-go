package main

import (
	"lucky-go/cloud"
	"lucky-go/finance"
	"lucky-go/forex"
	"lucky-go/game"
	"lucky-go/notify"
	"lucky-go/server/ssh"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd 表示在不带任何子命令的情况下调用时的基础命令
var rootCmd = &cobra.Command{
	Use:   "lucky-go",
	Short: "一个用于各种实用工具的CLI，包括云、金融、游戏和服务器操作",
	Long: `lucky-go 是一个CLI应用程序，提供与云服务交互、
分析金融数据、管理游戏自动化和服务器操作的实用工具。`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute 执行根命令并通过退出状态1处理任何错误。
// 此函数由 main.main() 调用，且只需执行一次。
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// 在这里定义标志和配置设置。
	// Cobra 支持持久标志，如果在此处定义，
	// 将对应用程序全局有效。

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lucky-go.yaml)")

	// Cobra 还支持本地标志，仅在直接调用此操作时运行。
	rootCmd.Flags().BoolP("toggle", "t", false, "切换选项的帮助消息")

	rootCmd.AddCommand(ssh.NewCommand())
	rootCmd.AddCommand(cloud.NewCommand())
	rootCmd.AddCommand(game.NewCommand())
	rootCmd.AddCommand(finance.NewCommand())
	rootCmd.AddCommand(notify.NewCommand())
	rootCmd.AddCommand(forex.NewCommand())
}
