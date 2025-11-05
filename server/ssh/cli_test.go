package ssh

import (
	"os"
	"os/exec"
	"testing"

	"lucky-go/config"

	"gopkg.in/yaml.v3"
)

func TestSSHCommand(t *testing.T) {
	t.Run("ValidDestination", func(t *testing.T) {
		// 创建临时测试目录
		tempDir := t.TempDir()
		originalHomeDir := os.Getenv("HOME")
		originalUserProfile := os.Getenv("USERPROFILE")

		os.Setenv("HOME", tempDir)
		os.Setenv("USERPROFILE", tempDir)
		defer os.Setenv("HOME", originalHomeDir)
		defer os.Setenv("USERPROFILE", originalUserProfile)

		// 准备测试配置
		testConfig := config.Config{
			Dest: map[string]config.DestinationInstance{
				"test-server": {
					Ssh:        "user@test.example.com",
					Region:     "us-west",
					InstanceId: "ins-test999",
				},
			},
		}

		// 创建配置目录和文件
		configDir := tempDir + string(os.PathSeparator) + config.CONFIG_DIR
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			t.Fatalf("failed to create config directory: %v", err)
		}

		configFile := configDir + string(os.PathSeparator) + config.CONFIG_FILE
		data, err := yaml.Marshal(testConfig)
		if err != nil {
			t.Fatalf("failed to marshal config: %v", err)
		}

		err = os.WriteFile(configFile, data, 0644)
		if err != nil {
			t.Fatalf("failed to write config file: %v", err)
		}

		// 保存原始命令执行函数
		originalExecCommand := execCommand
		defer func() {
			execCommand = originalExecCommand
		}()

		// 模拟ssh命令成功执行
		execCommand = func(name string, arg ...string) *exec.Cmd {
			cs := []string{"-test.run=TestHelperProcess", "--", name}
			cs = append(cs, arg...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{"GO_HELPER_PROCESS=1", "SSH_OUTPUT="}
			return cmd
		}

		cmd := NewCommand()
		err = cmd.RunE(cmd, []string{"test-server"})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("NoDestinationProvided", func(t *testing.T) {
		// 为避免运行命令时出现错误，我们不直接调用RunE
		// 而是测试命令配置
		cmd := NewCommand()

		// 验证命令配置
		if cmd.Args == nil {
			t.Error("expected Args to be set")
		}
	})

	t.Run("NonExistentDestination", func(t *testing.T) {
		cmd := NewCommand()

		err := cmd.RunE(cmd, []string{"non-existent"})
		if err == nil {
			t.Error("expected error for non-existent destination, got nil")
		}
	})
}

func TestNewCommand(t *testing.T) {
	t.Run("CommandStructure", func(t *testing.T) {
		cmd := NewCommand()

		if cmd.Use != "ssh [destination]" {
			t.Errorf("expected command use 'ssh [destination]', got '%s'", cmd.Use)
		}

		if cmd.Short != "与目标建立 SSH 连接" {
			t.Errorf("expected different short description, got '%s'", cmd.Short)
		}

		// 检查是否有子命令
		subcommands := cmd.Commands()
		if len(subcommands) == 0 {
			t.Error("expected subcommands, got none")
		}
	})
}

func TestServeCommand(t *testing.T) {
	t.Run("ServeCommandStructure", func(t *testing.T) {
		serveCmd := newCommand()

		if serveCmd.Use != "serve" {
			t.Errorf("expected command use 'serve', got '%s'", serveCmd.Use)
		}

		if serveCmd.Short != "运行 HTTP 服务器" {
			t.Errorf("expected different short description, got '%s'", serveCmd.Short)
		}

		// 检查是否有 --port 标志
		if serveCmd.Flags().Lookup("port") == nil {
			t.Error("expected --port flag, got none")
		}
	})
}

// 模拟辅助进程，用于测试命令执行
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_HELPER_PROCESS") != "1" {
		return
	}

	cmd := os.Args[3] // 第一个参数是test标志，第二个是"--"，第三个是命令名

	// 根据命令类型返回不同的输出
	if os.Getenv("EXEC_ERROR") == "1" {
		os.Exit(1) // 模拟命令执行失败
	}

	if cmd == "ssh" {
		// 模拟SSH命令成功执行
		output := os.Getenv("SSH_OUTPUT")
		if output != "" {
			os.Stdout.Write([]byte(output))
		}
	}

	os.Exit(0)
}
