package ssh

import (
	"errors"
	"fmt"
	"lucky-go/config"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh [destination]",
	Short: "build ssh connection with destination",
	Long:  `Connect to a remote server via SSH.\n\nYou must provide a destination in the format user@host.`,
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

		ssh := destinationInstance.Ssh

		fmt.Printf("get dest ssh %v\n", ssh)
		sshProcess := exec.Command("ssh", ssh)
		sshProcess.Stdin = os.Stdin
		sshProcess.Stdout = os.Stdout
		sshProcess.Stderr = os.Stderr

		if err := sshProcess.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to run ssh command, err: %v\n", err)
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
