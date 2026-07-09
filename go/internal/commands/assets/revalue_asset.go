package assets

import (
	"fmt"
	"strconv"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newRevalueAssetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "revalue-asset asset_id estimated_value",
		Short: "Change estimated asset value.",
		Args:  cobra.ExactArgs(2),
		RunE:  runRevalueAsset,
	}
}

func runRevalueAsset(cmd *cobra.Command, args []string) error {
	assetID := args[0]
	estimatedValue, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("estimated value must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.RevalueAsset(w, assetID, estimatedValue); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "revalue-asset", args); err != nil {
		return err
	}

	asset := w.Assets[assetID]
	commandlog.Action("Revalued asset: %s", assetID)
	commandlog.State("Owner: %s", asset.OwnerHumanID)
	commandlog.State("Estimated value: %d %s", asset.EstimatedValue, asset.Currency)
	return nil
}
