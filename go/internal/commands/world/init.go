package worldcmd

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create an empty banking world.",
		Args:  cobra.NoArgs,
		RunE:  runInit,
	}
}

func runInit(cmd *cobra.Command, args []string) error {
	alreadyExists := world.Exists()

	w := world.New()
	if err := commandrun.SaveWithHistory(w, "init", nil); err != nil {
		return err
	}

	if alreadyExists {
		commandlog.Action("Reset world: %s", world.FileName())
		return nil
	}

	commandlog.Action("Initialized empty world: %s", world.FileName())
	return nil
}
