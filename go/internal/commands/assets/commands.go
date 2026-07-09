package assets

import "github.com/spf13/cobra"

func Commands() []*cobra.Command {
	return []*cobra.Command{
		newRegisterAssetCommand(),
		newBuyAssetCashCommand(),
		newRevalueAssetCommand(),
	}
}
