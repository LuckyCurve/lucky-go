// Package ssh provides SSH connection utilities for the lucky-go application.
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
	Short: "Build SSH connection with destination",
	Long:  `Connect to a remote server via SSH using a destination name specified in the config.`,
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

// newCommand creates a subcommand for the ssh command that runs an HTTP server.
// It allows specifying a port to listen on with the --port flag.
func newCommand() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the HTTP server",
		Long:  `Run an HTTP server on the specified port.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Serving on :%d\n", port)
			return nil
		},
	}
	cmd.Flags().IntVar(&port, "port", 8080, "port to listen on")
	return cmd
}

// NewCommand creates and returns the SSH command with its subcommands for the server module.
func NewCommand() *cobra.Command {
	sshCmd.AddCommand(newCommand())

	return sshCmd
}
