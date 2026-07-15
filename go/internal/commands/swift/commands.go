package swift

import "github.com/spf13/cobra"

func Commands() []*cobra.Command {
	return []*cobra.Command{
		newOpenCorrespondentAccountCommand(),
		newFundCorrespondentAccountCommand(),
		newSwiftMT103Command(),
		newSettleSwiftCommand(),
		newRejectSwiftCommand(),
	}
}
