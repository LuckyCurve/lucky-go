package main

import (
	"lucky-go/cloud"
	"lucky-go/finance"
	"lucky-go/game"
	"lucky-go/server/ssh"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lucky-go",
	Short: "A CLI tool for various utilities including cloud, finance, game and server operations",
	Long: `lucky-go is a CLI application that provides utilities for interacting with cloud services, 
analyzing financial data, managing game automation, and server operations.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute executes the root command and handles any errors by exiting with status 1.
// This function is called by main.main() and only needs to be executed once.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lucky-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(ssh.NewCommand())
	rootCmd.AddCommand(cloud.NewCommand())
	rootCmd.AddCommand(game.NewCommand())
	rootCmd.AddCommand(finance.NewCommand())
}
