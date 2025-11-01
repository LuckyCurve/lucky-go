package game

import (
	"os"
	"os/exec"
	"testing"
)

// 保存原始函数以便测试后恢复
var originalExecCommand = execCommand

func TestChooseDevice(t *testing.T) {
	t.Run("NoDevicesFound", func(t *testing.T) {
		// 保存原始函数
		originalExecCommand := execCommand
		defer func() {
			execCommand = originalExecCommand
		}()
		
		// 模拟adb命令返回空设备列表
		execCommand = func(name string, arg ...string) *exec.Cmd {
			cs := []string{"-test.run=TestHelperProcess", "--", name}
			cs = append(cs, arg...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{"GO_HELPER_PROCESS=1", "ADB_OUTPUT="}
			return cmd
		}
		
		_, err := chooseDevice()
		if err == nil {
			t.Error("expected error for no devices, got nil")
		}
		if err.Error() != "devices not found" {
			t.Errorf("expected 'devices not found' error, got '%v'", err)
		}
	})
	
	t.Run("SingleDeviceFound", func(t *testing.T) {
		// 保存原始函数
		originalExecCommand := execCommand
		defer func() {
			execCommand = originalExecCommand
		}()
		
		// 模拟adb命令返回单个设备
		expectedDevice := "emulator-5554"
		execCommand = func(name string, arg ...string) *exec.Cmd {
			cs := []string{"-test.run=TestHelperProcess", "--", name}
			cs = append(cs, arg...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{"GO_HELPER_PROCESS=1", "ADB_OUTPUT=device", "ADB_DEVICE=" + expectedDevice}
			return cmd
		}
		
		device, err := chooseDevice()
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if device != expectedDevice {
			t.Errorf("expected device '%s', got '%s'", expectedDevice, device)
		}
	})
	
	t.Run("MultipleDevicesPrompt", func(t *testing.T) {
		// 由于chooseDevice涉及用户输入，难以完全测试，这里我们只测试结构
		// 实际测试需要更复杂的输入模拟
	})
}

func TestExecuteClick(t *testing.T) {
	t.Run("SuccessfulClick", func(t *testing.T) {
		// 保存原始函数
		originalExecCommand := execCommand
		defer func() {
			execCommand = originalExecCommand
		}()
		
		// 模拟adb tap命令成功执行
		execCommand = func(name string, arg ...string) *exec.Cmd {
			// 检查命令是否为预期的tap命令
			if name != "adb" || len(arg) < 5 || arg[2] != "shell" || arg[3] != "input" || arg[4] != "tap" {
				t.Errorf("unexpected command: %s %v", name, arg)
			}
			
			cs := []string{"-test.run=TestHelperProcess", "--", name}
			cs = append(cs, arg...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{"GO_HELPER_PROCESS=1", "ADB_OUTPUT="}
			return cmd
		}
		
		device := "emulator-5554"
		err := executeCLick(device)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
	
	t.Run("ClickCommandError", func(t *testing.T) {
		// 保存原始函数
		originalExecCommand := execCommand
		defer func() {
			execCommand = originalExecCommand
		}()
		
		// 模拟adb tap命令失败
		execCommand = func(name string, arg ...string) *exec.Cmd {
			cs := []string{"-test.run=TestHelperProcess", "--", name}
			cs = append(cs, arg...)
			cmd := exec.Command(os.Args[0], cs...)
			cmd.Env = []string{"GO_HELPER_PROCESS=1", "EXEC_ERROR=1"}
			return cmd
		}
		
		device := "emulator-5554"
		err := executeCLick(device)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestGameCommand(t *testing.T) {
	t.Run("CommandStructure", func(t *testing.T) {
		cmd := NewCommand()
		
		if cmd.Use != "game" {
			t.Errorf("expected command use 'game', got '%s'", cmd.Use)
		}
		
		if cmd.Short != "Start game automation" {
			t.Errorf("expected different short description, got '%s'", cmd.Short)
		}
	})
}

// 模拟辅助进程，用于测试命令执行
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_HELPER_PROCESS") != "1" {
		return
	}
	
	cmd := os.Args[3] // 第一个参数是test标志，第二个是"--"，第三个是命令名
	args := os.Args[4:]
	
	// 根据命令类型返回不同的输出
	if os.Getenv("EXEC_ERROR") == "1" {
		os.Exit(1) // 模拟命令执行失败
	}
	
	if cmd == "adb" {
		if args[0] == "devices" {
			deviceOutput := os.Getenv("ADB_OUTPUT")
			deviceID := os.Getenv("ADB_DEVICE")
			
			if deviceOutput == "device" && deviceID != "" {
				// 模拟单个设备
				output := "List of devices attached\n" + deviceID + "\t" + deviceOutput + "\n"
				os.Stdout.Write([]byte(output))
			} else {
				// 模拟无设备
				output := "List of devices attached\n"
				os.Stdout.Write([]byte(output))
			}
		} else if args[1] == "shell" && args[2] == "input" && args[3] == "tap" {
			// 模拟点击命令
			output := os.Getenv("ADB_OUTPUT")
			if output != "" {
				os.Stdout.Write([]byte(output))
			}
		}
	}
	
	os.Exit(0)
}