package world

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const DefaultFileName = "mini_world.json"

func FileName() string {
	if fileName := os.Getenv("MINI_WORLD_FILE"); fileName != "" {
		return fileName
	}

	return DefaultFileName
}

func Exists() bool {
	_, err := os.Stat(FileName())
	return err == nil
}

func Load() (*World, error) {
	fileName := FileName()
	data, err := os.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("world does not exist: %s. Run: go run ./cmd/mini init", fileName)
		}
		return nil, err
	}

	var world World

	if err := json.Unmarshal(data, &world); err != nil {
		return nil, err
	}

	return &world, nil
}

func Save(world *World) error {
	data, err := json.MarshalIndent(world, "", "  ")
	if err != nil {
		return err
	}

	data = append(data, '\n')
	return os.WriteFile(FileName(), data, 0o644)
}
