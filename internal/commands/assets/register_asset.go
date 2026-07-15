package assets

import (
	"fmt"
	"strconv"
	"strings"

	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newRegisterAssetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "register-asset asset_id owner_id currency estimated_value",
		Short: "Register an already existing asset.",
		Args:  cobra.ExactArgs(4),
		RunE:  runRegisterAsset,
	}
}

func runRegisterAsset(cmd *cobra.Command, args []string) error {
	assetID := args[0]
	ownerID := args[1]
	currency := strings.ToUpper(args[2])
	estimatedValue, err := strconv.Atoi(args[3])
	if err != nil {
		return fmt.Errorf("estimated value must be an integer")
	}

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.RegisterAsset(w, assetID, ownerID, currency, estimatedValue); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "register-asset", args); err != nil {
		return err
	}

	commandlog.Action("Registered asset: %s", assetID)
	asset := w.Assets[assetID]
	commandlog.State("Owner: %s %s", asset.OwnerType, asset.OwnerID)
	commandlog.State("Currency: %s", currency)
	commandlog.State("Estimated value: %d %s", estimatedValue, currency)
	return nil
}
