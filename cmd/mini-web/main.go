package main

import (
	"fmt"
	"os"

	"mini-world-go/internal/web"
)

func main() {
	addr := os.Getenv("MINI_WEB_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	if err := web.Run(addr); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
