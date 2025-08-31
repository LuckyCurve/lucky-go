package game

import (
	"fmt"
	"os/exec"
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

func executeCLick() error {
	cmd := exec.Command("adb", "-s", "127.0.0.1:5555", "shell", "input", "tap", "1800", "900")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if len(out) > 0 {
		fmt.Printf("execute command return: %v\n", string(out))
	}
	return nil
}

// func newCommand() *cobra.Command {
// 	var port int
// 	cmd := &cobra.Command{
// 		Use:   "serve",
// 		Short: "Run the HTTP server",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Printf("Serving on :%d\n", port)
// 			return nil
// 		},
// 	}
// 	cmd.Flags().IntVar(&port, "port", 8080, "port to listen on")
// 	return cmd
// }

func NewCommand() *cobra.Command {
	return sshCmd
}
