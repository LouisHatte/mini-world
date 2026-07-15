package swift

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newFundCorrespondentAccountCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "fund-correspondent-account correspondent_account_id amount",
		Short: "Fund a correspondent account using reserves.",
		Args:  cobra.ExactArgs(2),
		RunE:  runFundCorrespondentAccount,
	}
}

func runFundCorrespondentAccount(cmd *cobra.Command, args []string) error {
	correspondentAccountID := args[0]
	amount, err := parseAmount(args[1])
	if err != nil {
		return err
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	result, err := domain.FundCorrespondentAccount(w, correspondentAccountID, amount)
	if err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "fund-correspondent-account", args); err != nil {
		return err
	}

	account := w.CorrespondentAccounts[correspondentAccountID]
	centralBank := w.CentralBanks[result.CentralBankID]
	ownerBank := w.Banks[account.OwnerBankID]
	correspondentBank := w.Banks[account.CorrespondentBankID]
	commandlog.Action("Funded correspondent account %s with %d %s", correspondentAccountID, amount, account.Currency)
	logCorrespondentAccount(w, correspondentAccountID)
	commandlog.State("%s reserve account for %s: %d %s", result.CentralBankID, account.OwnerBankID, centralBank.ReserveAccounts[account.OwnerBankID], account.Currency)
	commandlog.State("%s reserve account for %s: %d %s", result.CentralBankID, account.CorrespondentBankID, centralBank.ReserveAccounts[account.CorrespondentBankID], account.Currency)
	commandlog.State("%s reserves at %s: %d %s", account.OwnerBankID, result.CentralBankID, ownerBank.ReserveBalances[result.CentralBankID], account.Currency)
	commandlog.State("%s reserves at %s: %d %s", account.CorrespondentBankID, result.CentralBankID, correspondentBank.ReserveBalances[result.CentralBankID], account.Currency)
	return nil
}
