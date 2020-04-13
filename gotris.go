package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't initialize screen and I'm sad.")
		os.Exit(-1)
	}

	g := NewGame(screen)
	g.Run()

	os.Exit(0)
}
