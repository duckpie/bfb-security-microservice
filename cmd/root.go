package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "server",
	}

	rootCmd.AddCommand(
		runCmd(),
	)

	return rootCmd
}
