package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't initialize screen and I'm sad.")
		os.Exit(-1)
	}

	rand.Seed(time.Now().UnixNano())
	g := NewGame(screen)
	g.Run()

	os.Exit(0)
}
