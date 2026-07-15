package worldcmd

import (
	"mini-world-go/internal/commandlog"
	"mini-world-go/internal/domain"
	"mini-world-go/internal/world"

	"github.com/spf13/cobra"
)

func newCheckWorldCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check-world",
		Short: "Check consistency between central banks and commercial banks.",
		Args:  cobra.NoArgs,
		RunE:  runCheckWorld,
	}
}

func runCheckWorld(cmd *cobra.Command, args []string) error {
	w, err := world.Load()
	if err != nil {
		return err
	}

	errors := domain.CheckWorld(w)

	if len(errors) == 0 {
		commandlog.Action("World check passed")
		commandlog.State("Reserve mirrors are consistent.")
		commandlog.State("Loan mirrors are consistent.")
		commandlog.State("Customer loan mirrors are consistent.")
		return nil
	}

	commandlog.Action("World check failed")

	for _, checkError := range errors {
		commandlog.State("- %s", checkError)
	}

	return nil
}
