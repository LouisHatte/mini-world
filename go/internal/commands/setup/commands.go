package setup

import "github.com/spf13/cobra"

func Commands() []*cobra.Command {
	return []*cobra.Command{
		newCreateCentralBankCommand(),
		newCreateBankCommand(),
		newCreateHumanCommand(),
		newOpenAccountCommand(),
	}
}
