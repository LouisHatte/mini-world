package commands

import (
	"sort"

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

type CommandInfo struct {
	Name  string `json:"name"`
	Use   string `json:"use"`
	Short string `json:"short"`
}

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
	return ExecuteArgs(nil)
}

func ExecuteArgs(args []string) error {
	rootCmd := newRootCommand()
	if args != nil {
		rootCmd.SetArgs(args)
	}
	return rootCmd.Execute()
}

func ListCommands() []CommandInfo {
	rootCmd := newRootCommand()
	infos := make([]CommandInfo, 0, len(rootCmd.Commands()))

	for _, command := range rootCmd.Commands() {
		if command.Hidden {
			continue
		}
		name := command.Name()
		if name == "completion" || name == "help" {
			continue
		}

		infos = append(infos, CommandInfo{
			Name:  name,
			Use:   command.Use,
			Short: command.Short,
		})
	}

	sort.Slice(infos, func(leftIndex, rightIndex int) bool {
		return infos[leftIndex].Name < infos[rightIndex].Name
	})

	return infos
}
