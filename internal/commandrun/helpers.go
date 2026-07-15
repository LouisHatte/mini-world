package commandrun

import "mini-world-go/internal/world"

func PrintBusinessError(err error) error {
	return err
}

func SaveWithHistory(w *world.World, command string, args []string) error {
	world.AppendCommandHistory(w, append([]string{command}, args...))
	return world.Save(w)
}
