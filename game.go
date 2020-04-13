package main

import (
	"time"
	"unicode"

	"github.com/gdamore/tcell"
)

type Game struct {
	screen   tcell.Screen
	board    Board
	shutdown chan bool
}

func NewGame(screen tcell.Screen) Game {
	g := new(Game)
	g.screen = screen
	g.board = *NewBoard()
	g.shutdown = make(chan bool)

	return *g
}

func (game *Game) Run() {
	initScreen(game.screen)
	defer game.screen.Fini()
	go game.getInput()

	for {
		select {
		case <-game.shutdown:
			return
		case <-time.After(time.Millisecond * 10):
		}
	}
}

func (game *Game) getInput() {
	for {
		event := game.screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventInterrupt:
			game.shutdown <- true
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape:
				game.shutdown <- true
			case tcell.KeyRune:
				if unicode.ToLower(event.Rune()) == 'q' {
					game.shutdown <- true
				}
			}
		}
	}
}
