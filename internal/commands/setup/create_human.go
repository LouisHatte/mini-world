package setup

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newCreateHumanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create-human human_id",
		Short: "Create a human actor.",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreateHuman,
	}
}

func runCreateHuman(cmd *cobra.Command, args []string) error {
	humanID := args[0]

	w, err := world.Load()
	if err != nil {
		return err
	}

	if err := domain.CreateHuman(w, humanID); err != nil {
		return commandrun.PrintBusinessError(err)
	}

	if err := commandrun.SaveWithHistory(w, "create-human", args); err != nil {
		return err
	}

	commandlog.Action("Created human: %s", humanID)
	return nil
}
