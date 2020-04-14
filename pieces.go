package main

import "math/rand"

type Point struct{ X, Y int }

func (point *Point) Add(offset Point) Point {
	return Point{point.X + offset.X, point.Y + offset.Y}
}

type Piece [4]Point

// --X-
// --X-
// --X-
// --X-
var pieceLongBoi = Piece{Point{0, -2}, Point{0, -1}, Point{0, 0}, Point{0, 1}}

// ----  ----
// -X--  --X-
// -X--  --X-
// -XX-  -XX-
var pieceL = Piece{Point{-1, -1}, Point{-1, 0}, Point{-1, 1}, Point{0, 1}}
var pieceInvertedL = Piece{Point{0, -1}, Point{0, 0}, Point{-1, 1}, Point{0, 1}}

// ----  ----
// -X--  --X-
// -XX-  -XX-
// --X-  -X--
var pieceZ = Piece{Point{-1, -1}, Point{-1, 0}, Point{0, 0}, Point{0, 1}}
var pieceInvertedZ = Piece{Point{0, -1}, Point{-1, 0}, Point{0, 0}, Point{-1, 1}}

// ----
// -XX-
// -XX-
// ----
var pieceBlock = Piece{Point{-1, -1}, Point{0, -1}, Point{-1, 0}, Point{0, 0}}

// ----
// -X--
// XXX-
// ----
var pieceTBoi = Piece{Point{0, -1}, Point{-1, 0}, Point{0, 0}, Point{1, 0}}

func (piece *Piece) RotatePiece(clockwise bool) *Piece {
	rotated := new(Piece)
	for i := 0; i < len(piece); i++ {
		if clockwise == true {
			rotated[i] = Point{piece[i].Y, -piece[i].X}
		} else {
			rotated[i] = Point{-piece[i].Y, piece[i].X}
		}
	}
	return rotated
}

type BagOfPieces struct {
	pieces []Piece
	order  []int
	next   int
}

func NewBagOfPieces() *BagOfPieces {
	pieces := []Piece{pieceLongBoi, pieceL, pieceInvertedL, pieceZ, pieceInvertedZ, pieceBlock, pieceTBoi}
	order := []int{0, 1, 2, 3, 4, 5, 6}

	return &BagOfPieces{pieces: pieces, order: order, next: len(order)}
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

	return &bag.pieces[bag.order[bag.next]]
}

func (bag *BagOfPieces) TakeNextPiece() *Piece {
	ret := bag.NextPiece()
	bag.next++
	return ret
}
