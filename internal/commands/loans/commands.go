package loans

import "github.com/spf13/cobra"

func Commands() []*cobra.Command {
	return []*cobra.Command{
		newGrantLoanCommand(),
		newAccrueInterestCommand(),
		newRepayLoanCommand(),
		newDefaultLoanCommand(),
	}
}
