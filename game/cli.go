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

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "game",
	Short: "start agme hang out",
	RunE: func(cmd *cobra.Command, args []string) error {

		for {
			err := executeCLick()
			if err != nil {
				return err
			}

			time.Sleep(5 * time.Second)
		}
	},
}

func chooseDevice() (string, error) {
	cmd := exec.Command("adb", "devices")
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
		return "", fmt.Errorf("devices not found!")
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

func executeCLick() error {
	device, err := chooseDevice()
	if err != nil {
		return err
	}

	cmd := exec.Command("adb", "-s", device, "shell", "input", "tap", "1800", "900")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if len(out) > 0 {
		fmt.Printf("execute command return: %v\n", string(out))
	}
	return nil
}

func NewCommand() *cobra.Command {
	return sshCmd
}
