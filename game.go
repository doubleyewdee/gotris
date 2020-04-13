package main

import (
	"fmt"
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
	game.screen.Init()
	game.screen.Clear()
	defer game.screen.Fini()

	go game.getInput()

	drawCount := 0
	for {
		game.writeMsg("awaiting input...")
		select {
		case <-game.shutdown:
			return
		case <-time.After(time.Millisecond * 500):
		}
		game.writeMsg("drawing...")

		drawCount++
		for y := 0; y < BOARD_HEIGHT; y++ {
			for x := 0; x < BOARD_WIDTH; x++ {
				if (x+y+drawCount)%2 == 0 {
					game.board.Cells[y][x].Color = tcell.NewRGBColor(255, 255, 255)
				} else {
					game.board.Cells[y][x].Color = tcell.NewRGBColor(0, 0, 0)
				}
			}
		}
		game.writeMsg(fmt.Sprintf("dc:%v", drawCount))

		game.draw()
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

const BLOCK_CHAR = '\u2592' //

func (game *Game) draw() {
	// XXX: figure out how to center the board lol
	charStyle := tcell.StyleDefault
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			charStyle = charStyle.Background(game.board.Cells[y][x].Color)
			game.screen.SetContent(x*2, y, BLOCK_CHAR, nil, charStyle)
			game.screen.SetContent((x*2)+1, y, BLOCK_CHAR, nil, charStyle)
		}
	}
	game.screen.Show()
}

func (game *Game) writeMsg(str string) {
	y := BOARD_HEIGHT + 1

	width, _ := game.screen.Size()
	for x := 0; x < width; x++ {
		game.screen.SetContent(x, y, '\u0000', nil, tcell.StyleDefault)
	}

	for x, ch := range str {
		game.screen.SetContent(x, y, ch, nil, tcell.StyleDefault)
	}
	game.screen.Show()
}
