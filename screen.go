package main

import (
	"github.com/gdamore/tcell"
)

type Color struct {
	R int
	G int
	B int
}

func initScreen(screen tcell.Screen) {
	screen.Init()
	screen.Clear()
	return
}
