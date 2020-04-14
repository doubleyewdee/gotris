package main

import (
	"github.com/gdamore/tcell"
)

const (
	BOARD_WIDTH  = 10
	BOARD_HEIGHT = 20
)

const (
	PIECE_VALID = iota
	PIECE_INVALID
	PIECE_OVERLAPPED
)

type Cell struct {
	Color  tcell.Color
	Locked bool
}

type Board struct {
	Cells [][]Cell
}

func NewBoard() *Board {
	cells := make([][]Cell, BOARD_HEIGHT)
	for i := range cells {
		cells[i] = make([]Cell, BOARD_WIDTH)
	}

	return &Board{Cells: cells}
}

// attempt to place the piece at the given position, will return false
// if the addition would be invalid
func (board *Board) PlacePiece(piece Piece, position Point) int {
	if valid, _ := board.isPiecePositionValid(piece, position); !valid {
		return PIECE_INVALID
	} else if board.isPiecePositionOverlapped(piece, position) {
		return PIECE_OVERLAPPED
	} else {
		return PIECE_VALID
	}
}

// locks all intersecting spots on the board for the given piece
func (board *Board) LockPiece(piece Piece, position Point) {
	topY, bottomY := BOARD_HEIGHT, 0

	r, g, b := piece.color.RGB()
	r -= r / 4
	g -= g / 4
	b -= b / 4
	cellColor := tcell.NewRGBColor(r, g, b)
	locker := func(point Point) {
		cell := board.cellAt(point)
		cell.Locked = true
		cell.Color = cellColor
		if point.Y < topY {
			topY = point.Y
		}
		if point.Y > bottomY {
			bottomY = point.Y
		}
	}

	board.intersectPiece(piece, position, locker)

	sweepLine := bottomY
	for sweepLine >= topY {
		cleared := true
		for _, cell := range board.Cells[sweepLine] {
			if !cell.Locked {
				cleared = false
				break
			}
		}

		if cleared {
			board.ClearLine(sweepLine)
			topY--
		} else {
			sweepLine--
		}
	}
}

// clears a line (woo!) and shifts all the lines above down
func (board *Board) ClearLine(line int) {
	for l := line; l > 0; l-- {
		board.Cells[l] = board.Cells[l-1]
	}
	board.Cells[0] = make([]Cell, BOARD_WIDTH)
}

// returns whether the piece's position is valid (within the board's bounds), and if not, what change to position would
// make it fit
func (board *Board) isPiecePositionValid(piece Piece, position Point) (bool, Point) {
	valid := true
	minX, minY := 0, 0
	maxX, maxY := BOARD_WIDTH-1, BOARD_HEIGHT-1

	validator := func(point Point) {
		// if we can't get the cell try and figure out how far out of bounds it is
		cell := board.cellAt(point)
		if cell == nil {
			valid = false
			if point.X < minX {
				minX = point.X
			}
			if point.Y < minY {
				minY = point.Y
			}
			if point.X > maxX {
				maxX = point.X
			}
			if point.Y > maxY {
				maxY = point.Y
			}
		}
	}

	board.intersectPiece(piece, position, validator)

	offset := Point{}
	if minX < 0 {
		offset.X = -minX
	} else if maxX >= BOARD_WIDTH {
		offset.X = -(maxX - (BOARD_WIDTH - 1))
	}

	if minY < 0 {
		offset.Y = -minY
	} else if maxY >= BOARD_HEIGHT {
		offset.Y = -(maxY - (BOARD_HEIGHT - 1))
	}

	return valid, offset
}

func (board *Board) isPiecePositionOverlapped(piece Piece, position Point) bool {
	overlapped := false
	validator := func(point Point) {
		cell := board.cellAt(point)
		if cell != nil && cell.Locked {
			overlapped = true
		}
	}

	board.intersectPiece(piece, position, validator)
	return overlapped
}

// returns Points relative to the board for the current piece and a flag if one or more points are invalid
func (board *Board) intersectPiece(piece Piece, position Point, pointOperator func(Point)) {
	for i := 0; i < len(piece.points); i++ {
		pt := piece.points[i]
		x := position.X + pt.X
		y := position.Y + pt.Y
		pointOperator(Point{x, y})
	}
}

func (board *Board) cellAt(point Point) *Cell {
	if point.X < 0 || point.X >= BOARD_WIDTH || point.Y < 0 || point.Y >= BOARD_HEIGHT {
		return nil
	}

	return &board.Cells[point.Y][point.X]
}
