package cloud

import (
	"errors"
	"fmt"
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

func newCommand() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Serving on :%d\n", port)
			return nil
		},
	}
	cmd.Flags().IntVar(&port, "port", 8080, "port to listen on")
	return cmd
}

func NewCommand() *cobra.Command {
	sshCmd.AddCommand(newCommand())

	return sshCmd
}
