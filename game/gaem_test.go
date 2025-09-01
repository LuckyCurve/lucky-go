package game

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestXxx(t *testing.T) {
	// 执行 adb devices
	cmd := exec.Command("adb", "devices")
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	// 转成字符串
	result := string(out)
	lines := strings.Split(result, "\n")

	var devices []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 过滤掉空行和头部 "List of devices attached"
		if line == "" || strings.HasPrefix(line, "List of devices") {
			continue
		}
		// adb devices 输出的格式一般是 "序列号 \t 状态"
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "device" {
			devices = append(devices, fields[0])
		}
	}

	// 打印结果
	fmt.Println("Devices:", devices)
}
