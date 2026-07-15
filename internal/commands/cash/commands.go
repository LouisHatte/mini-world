package cash

import "github.com/spf13/cobra"

func Commands() []*cobra.Command {
	return []*cobra.Command{
		newIssueCashCommand(),
		newSeedCashCommand(),
		newTransferCashCommand(),
		newDepositCashCommand(),
		newWithdrawCashCommand(),
		newSupplyCashCommand(),
		newMoveCashCommand(),
		newSellCashCommand(),
		newReturnCashCommand(),
		newDestroyCashCommand(),
	}
}
