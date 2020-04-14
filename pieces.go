package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

type Point struct{ X, Y int }

func (point *Point) Add(offset Point) Point {
	return Point{point.X + offset.X, point.Y + offset.Y}
}

const (
	PIECE_I = iota
	PIECE_J
	PIECE_L
	PIECE_O
	PIECE_S
	PIECE_T
	PIECE_Z
)

var DEFAULT_PIECE_COLOR = tcell.NewRGBColor(128, 48, 192)

type Piece struct {
	points [4]Point
	color  tcell.Color
	id     int
}

// --X-
// --X-
// --X-
// --X-
var pieceI = Piece{points: [4]Point{{0, -2}, {0, -1}, {0, 0}, {0, 1}}, id: PIECE_I}

// ----  ----
// -X--  --X-
// -X--  --X-
// -XX-  -XX-
var pieceL = Piece{points: [4]Point{{-1, -1}, {-1, 0}, {-1, 1}, {0, 1}}, id: PIECE_L}
var pieceJ = Piece{points: [4]Point{{0, -1}, {0, 0}, {-1, 1}, {0, 1}}, id: PIECE_J}

// ----  ----
// -X--  --X-
// -XX-  -XX-
// --X-  -X--
var pieceS = Piece{points: [4]Point{{-1, -1}, {-1, 0}, {0, 0}, {0, 1}}, id: PIECE_S}
var pieceZ = Piece{points: [4]Point{{0, -1}, {-1, 0}, {0, 0}, {-1, 1}}, id: PIECE_Z}

// ----
// -XX-
// -XX-
// ----
var pieceO = Piece{points: [4]Point{{-1, -1}, {0, -1}, {-1, 0}, {0, 0}}, id: PIECE_O}

// ----
// -X--
// XXX-
// ----
var pieceT = Piece{points: [4]Point{{0, -1}, {-1, 0}, {0, 0}, {1, 0}}, id: PIECE_T}

func (piece *Piece) RotatePiece(clockwise bool) *Piece {
	rotated := new([4]Point)
	for i := 0; i < len(piece.points); i++ {
		if clockwise == true {
			rotated[i] = Point{piece.points[i].Y, -piece.points[i].X}
		} else {
			rotated[i] = Point{-piece.points[i].Y, piece.points[i].X}
		}
	}
	return &Piece{points: *rotated, id: piece.id, color: piece.color}
}

type BagOfPieces struct {
	pieces   []Piece
	order    []int
	next     int
	colorMap map[int]tcell.Color
}

func NewBagOfPieces() *BagOfPieces {
	pieces := []Piece{pieceI, pieceJ, pieceL, pieceO, pieceS, pieceT, pieceZ}
	order := []int{PIECE_I, PIECE_J, PIECE_L, PIECE_O, PIECE_S, PIECE_T, PIECE_Z}
	colorMap := make(map[int]tcell.Color, len(pieces))
	for _, p := range pieces {
		colorMap[p.id] = DEFAULT_PIECE_COLOR
	}

	return &BagOfPieces{pieces: pieces, order: order, next: len(order), colorMap: colorMap}
}

func (bag *BagOfPieces) NextPiece() *Piece {
	if bag.next == len(bag.order) {
		// re-shuffle
		for i := len(bag.order) - 1; i > 0; i-- {
			j := rand.Intn(i)
			if j != i {
				bag.order[i], bag.order[j] = bag.order[j], bag.order[i]
			}
		}

		bag.next = 0
	}

	piece := bag.pieces[bag.order[bag.next]]
	return &Piece{points: piece.points, id: piece.id, color: bag.colorMap[piece.id]}
}

func (bag *BagOfPieces) TakeNextPiece() *Piece {
	ret := bag.NextPiece()
	bag.next++
	return ret
}
