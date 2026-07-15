package payments

import "github.com/spf13/cobra"

func Commands() []*cobra.Command {
	return []*cobra.Command{
		newInternalTransferCommand(),
		newInterbankPaymentCommand(),
		newPayCommand(),
	}
}
