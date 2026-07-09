package commandrun

import (
	"fmt"

	"mini-world-go/internal/world"
)

func PrintBusinessError(err error) error {
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func SaveWithHistory(w *world.World, command string, args []string) error {
	world.AppendCommandHistory(w, append([]string{command}, args...))
	return world.Save(w)
}
