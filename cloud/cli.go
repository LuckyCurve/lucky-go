package cloud

import (
	"errors"
	"lucky-go/config"

	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "reboot [destination]",
	Short: "reboot destination machine",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		destination := args[0]

		if destination == "" {
			return errors.New("destination must exists")
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

func NewCommand() *cobra.Command {
	return sshCmd
}
