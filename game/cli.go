// Package game 提供游戏自动化实用工具。
package game

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// gameCmd 表示游戏自动化命令
var gameCmd = &cobra.Command{
	Use:   "game",
	Short: "启动游戏自动化",
	Long:  `开始游戏挂机的自动点击。`,
	RunE: func(cmd *cobra.Command, args []string) error {

		device, err := chooseDevice()
		if err != nil {
			return err
		}

		for {
			err := executeCLick(device)
			if err != nil {
				return err
			}

			time.Sleep(5 * time.Second)
		}
	},
}

// chooseDevice 提示用户从连接的设备列表中选择一个Android设备。
// 如果未找到设备或用户输入无效，则返回所选设备ID或错误。
func chooseDevice() (string, error) {
	cmd := execCommand("adb", "devices")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	result := string(out)
	lines := strings.Split(result, "\n")

	var devices []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "List of devices") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "device" {
			devices = append(devices, fields[0])
		}
	}

	if len(devices) == 0 {
		return "", fmt.Errorf("未找到设备")
	} else if len(devices) == 1 {
		return devices[0], nil
	} else {
		fmt.Println("检测到多个设备，请选择:")
		for i, dev := range devices {
			fmt.Printf("[%d] %s\n", i, dev)
		}

		fmt.Print("请输入下标: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		idx, err := strconv.Atoi(input)
		if err != nil || idx < 0 || idx >= len(devices) {
			return "", fmt.Errorf("输入无效")
		}

		fmt.Println("你选择的设备是:", devices[idx])

		return devices[idx], nil
	}
}

// 为了测试目的，定义一个可替换的执行命令函数
var execCommand = exec.Command

// executeCLick 在指定设备的坐标(1800, 900)上执行ADB点击命令。
// 这用于游戏自动化以执行点击操作。
func executeCLick(device string) error {
	cmd := execCommand("adb", "-s", device, "shell", "input", "tap", "1800", "900")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if len(out) > 0 {
		fmt.Printf("执行命令返回: %v\n", string(out))
	}
	return nil
}

// NewCommand 为游戏模块创建并返回游戏自动化命令。
func NewCommand() *cobra.Command {
	return gameCmd
}
