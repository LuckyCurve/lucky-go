// Package ssh 为 lucky-go 应用程序提供 SSH 连接实用工具。
package ssh

import (
	"errors"
	"fmt"
	"lucky-go/config"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// 为测试目的定义一个可替换的执行命令函数
var execCommand = exec.Command

// sshCmd 表示 ssh 命令
var sshCmd = &cobra.Command{
	Use:   "ssh [destination]",
	Short: "与目标建立 SSH 连接",
	Long:  `使用配置中指定的目标名称通过 SSH 连接到远程服务器。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("必须提供目标")
		}

		destination := args[0]

		if destination == "" {
			return errors.New("必须提供目标")
		}

		destinationInstance, err := config.LoadDestinationInstance(destination)
		if err != nil {
			return err
		}

		ssh := destinationInstance.Ssh

		fmt.Printf("get dest ssh %v\n", ssh)
		sshProcess := execCommand("ssh", ssh)
		sshProcess.Stdin = os.Stdin
		sshProcess.Stdout = os.Stdout
		sshProcess.Stderr = os.Stderr

		if err := sshProcess.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run ssh command, err: %v\n", err)
			return err
		}

		return nil
	},
}

// newCommand 为 ssh 命令创建一个运行 HTTP 服务器的子命令。
// 它允许使用 --port 标志指定要监听的端口。
func newCommand() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "运行 HTTP 服务器",
		Long:  `在指定端口上运行 HTTP 服务器。`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Serving on :%d\n", port)
			return nil
		},
	}
	cmd.Flags().IntVar(&port, "port", 8080, "要监听的端口")
	return cmd
}

// NewCommand 为服务器模块创建并返回带有子命令的 SSH 命令。
func NewCommand() *cobra.Command {
	sshCmd.AddCommand(newCommand())

	return sshCmd
}
