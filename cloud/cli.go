package cloud

import (
	"errors"
	"lucky-go/config"

	"github.com/spf13/cobra"
)

// rebootCmd represents the reboot command
var rebootCmd = &cobra.Command{
	Use:   "reboot [destination]",
	Short: "Reboot destination machine",
	Long:  `Reboot a cloud instance specified by the destination name.`,
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

// NewCommand creates and returns the reboot command for the cloud module.
func NewCommand() *cobra.Command {
	return rebootCmd
}
