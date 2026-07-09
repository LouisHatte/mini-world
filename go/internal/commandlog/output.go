package commandlog

import "fmt"

func Action(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func State(format string, args ...any) {
	fmt.Printf("\t"+format+"\n", args...)
}
