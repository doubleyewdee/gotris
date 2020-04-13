package main

import (
	"fmt"
	"time"
	"unicode"

	"github.com/gdamore/tcell"
)

const DRAWN_CELL_WIDTH = 2 // 2-width cells look vaguely like a square
var DRAW_COLOR = tcell.NewRGBColor(128, 48, 192)
var EMPTY_COLOR = tcell.NewRGBColor(0, 0, 0)

const (
	COMMAND_QUIT = iota
	COMMAND_NEW_GAME
	COMMAND_ROTATE
	COMMAND_ROTATE_COUNTER
	COMMAND_LEFT
	COMMAND_RIGHT
	COMMAND_DOWN
)

type Game struct {
	screen         tcell.Screen
	board          Board
	command        chan int
	boardArea      struct{ X, Y, Width, Height int }
	messageLine    int
	pieceGenerator BagOfPieces
	currentPiece   *Piece
	piecePosition  Point
	advanceSpeed   time.Duration
	lastAdvanced   time.Time
}

func NewGame(screen tcell.Screen) Game {
	g := new(Game)
	g.screen = screen
	g.board = *NewBoard()
	g.pieceGenerator = *NewBagOfPieces()
	g.command = make(chan int)
	g.boardArea.Height = BOARD_HEIGHT
	g.boardArea.Width = BOARD_WIDTH * DRAWN_CELL_WIDTH
	g.advanceSpeed = time.Second

	return *g
}

func (game *Game) Run() {
	game.screen.Init()
	game.screen.HideCursor()
	game.screen.Clear()
	defer game.screen.Fini()

	game.setLayout()

	drawCount := 0
	lastCommand := COMMAND_QUIT

	go game.getInput()
	for {
		now := time.Now()
		advanceTime := time.Second
		if game.currentPiece != nil {
			advanceTime = game.advanceSpeed - now.Sub(game.lastAdvanced)
			if advanceTime <= 0 {
				if !game.tryMovePiece(Point{0, 1}) {
					// trying to move a piece down again if you can't means setting the piece
					game.lockPiece()
				} else {
					game.lastAdvanced = now
				}
				continue
			}
		}

		select {
		case cmd := <-game.command:
			switch cmd {
			case COMMAND_QUIT:
				return
			case COMMAND_NEW_GAME:
				game.board = *NewBoard()
			case COMMAND_ROTATE:
				game.tryRotatePiece()
			case COMMAND_LEFT:
				game.tryMovePiece(Point{-1, 0})
			case COMMAND_RIGHT:
				game.tryMovePiece(Point{1, 0})
			case COMMAND_DOWN:
				if !game.tryMovePiece(Point{0, 1}) {
					// trying to move a piece down again if you can't means setting the piece
					game.lockPiece()
				} else {
					game.lastAdvanced = now
				}
			}
			lastCommand = cmd
		case <-time.After(advanceTime):
			game.writeMsg("lc:%v, dc:%v, piece:%v, pos:%v",
				lastCommand,
				drawCount,
				game.currentPiece,
				game.piecePosition)
			//game.command <- COMMAND_DOWN
		}

		if game.currentPiece == nil && !game.tryAddPiece() {
			game.writeMsg("Game over man, game over!")
			continue
		}

		drawCount++
		game.draw()
		game.screen.Show()
	}
}

func (game *Game) tryRotatePiece() bool {
	if game.currentPiece == nil {
		return false
	}

	piece := game.currentPiece.RotatePiece(true)

	ret := game.board.PlacePiece(*piece, game.piecePosition)
	switch ret {
	case PIECE_VALID:
		game.currentPiece = piece
	case PIECE_INVALID:
		_, shift := game.board.isPiecePositionValid(*piece, game.piecePosition)
		newPosition := Point{game.piecePosition.X + shift.X, game.piecePosition.Y + shift.Y}
		if game.board.PlacePiece(*piece, newPosition) == PIECE_VALID {
			game.piecePosition = newPosition
			game.currentPiece = piece
		}
	default:
		return false
	}
	return true
}

func (game *Game) tryAddPiece() bool {
	piece, position := new(Piece), Point{BOARD_WIDTH / 2, 0}
	piece = game.pieceGenerator.NextPiece()
	if valid, shift := game.board.isPiecePositionValid(*piece, position); !valid {
		position = Point{position.X + shift.X, position.Y + shift.Y}
	}
	if game.board.isPiecePositionOverlapped(*piece, position) {
		return false
	}

	game.currentPiece = piece
	game.piecePosition = position
	game.lastAdvanced = time.Now()
	return true
}

func (game *Game) tryMovePiece(offset Point) bool {
	if game.currentPiece == nil {
		return false
	}

	newPosition := game.piecePosition.Add(offset)
	if valid, _ := game.board.isPiecePositionValid(*game.currentPiece, newPosition); valid &&
		!game.board.isPiecePositionOverlapped(*game.currentPiece, newPosition) {
		game.piecePosition = newPosition
		return true
	}

	return false
}

func (game *Game) lockPiece() {
	if game.currentPiece == nil {
		return
	}

	game.board.LockPiece(*game.currentPiece, game.piecePosition)
	game.currentPiece = nil
}

func (game *Game) setLayout() {
	// center the board
	width, height := game.screen.Size()
	width = width/2 - (BOARD_WIDTH*DRAWN_CELL_WIDTH)/2
	height = height/2 - BOARD_HEIGHT/2
	game.boardArea.X, game.boardArea.Y = width, height

	game.messageLine = game.boardArea.Y - 2
}

func (game *Game) getInput() {
	for {
		event := game.screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventInterrupt:
			game.command <- COMMAND_QUIT
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape:
				game.command <- COMMAND_QUIT
			case tcell.KeyLeft:
				game.command <- COMMAND_LEFT
			case tcell.KeyRight:
				game.command <- COMMAND_RIGHT
			case tcell.KeyDown:
				game.command <- COMMAND_DOWN
			case tcell.KeyRune:
				switch unicode.ToLower(event.Rune()) {
				case 'n':
					game.command <- COMMAND_NEW_GAME
				case 'q':
					game.command <- COMMAND_QUIT
				case 'r':
					game.command <- COMMAND_ROTATE
				}
			}
		}
	}
}

const (
	BLOCK_CHAR = '\u2592'
	EMPTY_CHAR = '\u0000'
)

func (game *Game) draw() {
	charStyle := tcell.StyleDefault
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			charStyle = charStyle.Background(game.board.Cells[y][x].Color)
			ch := EMPTY_CHAR
			if game.board.Cells[y][x].Locked {
				ch = BLOCK_CHAR
			}
			for i := 0; i < DRAWN_CELL_WIDTH; i++ {
				game.screen.SetContent(
					game.boardArea.X+(x*DRAWN_CELL_WIDTH)+i,
					game.boardArea.Y+y,
					ch, nil, charStyle)
			}
		}
	}

	for i := 0; i < len(game.currentPiece); i++ {
		piece := game.currentPiece[i]
		x := game.piecePosition.X + piece.X
		y := game.piecePosition.Y + piece.Y
		charStyle = charStyle.Background(EMPTY_COLOR)
		charStyle = charStyle.Foreground(DRAW_COLOR)
		for i := 0; i < DRAWN_CELL_WIDTH; i++ {
			game.screen.SetContent(
				game.boardArea.X+(x*DRAWN_CELL_WIDTH)+i,
				game.boardArea.Y+y,
				BLOCK_CHAR, nil, charStyle)
		}
	}
}

func (game *Game) writeMsg(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)

	width, _ := game.screen.Size()
	st := tcell.StyleDefault.Reverse(true)
	for x := 0; x < width; x++ {
		game.screen.SetContent(x, game.messageLine, '\u0000', nil, st)
	}

	startX := width/2 - len(str)/2
	for x, ch := range str {
		game.screen.SetContent(startX+x, game.messageLine, ch, nil, st)
	}
	game.screen.Show()
}
