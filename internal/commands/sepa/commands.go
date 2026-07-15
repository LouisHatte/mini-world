package sepa

import "github.com/spf13/cobra"

func Commands() []*cobra.Command {
	return []*cobra.Command{
		newSepaCreditTransferCommand(),
		newSepaInstantCommand(),
		newSettleSepaCommand(),
		newRejectSepaCommand(),
	}
}
