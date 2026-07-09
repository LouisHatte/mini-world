package commands

import (
	"mini-world-go/internal/commands/cash"
	"mini-world-go/internal/commands/setup"
	worldcmd "mini-world-go/internal/commands/world"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mini",
		Short: "A small monetary and banking simulation CLI.",
	}

	rootCmd.AddCommand(worldcmd.Commands()...)
	rootCmd.AddCommand(setup.Commands()...)
	rootCmd.AddCommand(cash.Commands()...)

	return rootCmd
}

func Execute() error {
	return newRootCommand().Execute()
}
