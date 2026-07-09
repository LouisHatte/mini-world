package world

import "time"

func AppendCommandHistory(world *World, argv []string) {
	var command *string

	if len(argv) > 0 {
		command = &argv[0]
	}

	world.CommandHistory = append(world.CommandHistory, CommandHistoryEntry{
		ID:           len(world.CommandHistory) + 1,
		TimestampUTC: time.Now().UTC().Format(time.RFC3339Nano),
		Command:      command,
		Argv:         argv,
	})
}
