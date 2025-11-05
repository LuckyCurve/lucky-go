package cloud

import (
	"errors"
	"lucky-go/config"

	"github.com/spf13/cobra"
)

// rebootCmd 表示重启命令
var rebootCmd = &cobra.Command{
	Use:   "reboot [destination]",
	Short: "重启目标机器",
	Long:  `重启由目标名称指定的云实例。`,
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

		err = RebootInstance(destinationInstance)

		if err != nil {
			return err
		}

		return nil
	},
}

// NewCommand 为云模块创建并返回重启命令。
func NewCommand() *cobra.Command {
	return rebootCmd
}
