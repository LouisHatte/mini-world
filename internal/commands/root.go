package commands

import (
	"mini-world-go/internal/commands/assets"
	"mini-world-go/internal/commands/cash"
	"mini-world-go/internal/commands/fx"
	"mini-world-go/internal/commands/loans"
	"mini-world-go/internal/commands/payments"
	"mini-world-go/internal/commands/reserves"
	"mini-world-go/internal/commands/sepa"
	"mini-world-go/internal/commands/setup"
	"mini-world-go/internal/commands/swift"
	worldcmd "mini-world-go/internal/commands/world"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "mini",
		Short:         "A small monetary and banking simulation CLI.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(worldcmd.Commands()...)
	rootCmd.AddCommand(setup.Commands()...)
	rootCmd.AddCommand(cash.Commands()...)
	rootCmd.AddCommand(assets.Commands()...)
	rootCmd.AddCommand(reserves.Commands()...)
	rootCmd.AddCommand(payments.Commands()...)
	rootCmd.AddCommand(loans.Commands()...)
	rootCmd.AddCommand(sepa.Commands()...)
	rootCmd.AddCommand(fx.Commands()...)
	rootCmd.AddCommand(swift.Commands()...)

	return rootCmd
}

func Execute() error {
	return newRootCommand().Execute()
}
