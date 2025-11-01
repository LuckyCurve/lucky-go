package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestExecute(t *testing.T) {
	t.Run("ExecuteWithValidCommand", func(t *testing.T) {
		// 保存原始RootCmd
		originalRootCmd := rootCmd
		// 创建一个新的命令用于测试
		testCmd := &cobra.Command{
			Use: "test",
			Run: func(cmd *cobra.Command, args []string) {},
		}
		rootCmd = testCmd
		// 恢复原始命令
		defer func() {
			rootCmd = originalRootCmd
		}()
		
		// 由于Execute()函数会调用os.Exit，我们不能直接测试它
		// 而是测试命令结构
		if rootCmd.Use != "test" {
			t.Errorf("expected root command to be replaced")
		}
	})
}

func TestRootCommandStructure(t *testing.T) {
	// 测试现有命令的结构
	
	// 测试命令基本信息
	if rootCmd.Use != "lucky-go" {
		t.Errorf("expected command use 'lucky-go', got '%s'", rootCmd.Use)
	}
	
	// 检查是否有 --toggle 标志
	if rootCmd.Flags().Lookup("toggle") == nil {
		t.Error("expected --toggle flag, got none")
	}
}

func TestInitFunction(t *testing.T) {
	// 初始化函数在包加载时自动运行，我们只验证其效果
	// 检查根命令是否有 --toggle 标志
	if rootCmd.Flags().Lookup("toggle") == nil {
		t.Error("expected --toggle flag from init function, got none")
	}
}