package main

type Point struct{ X, Y int }

func (point *Point) Add(offset Point) Point {
	return Point{point.X + offset.X, point.Y + offset.Y}
}

type Piece [4]Point
type BagOfPieces struct {
	pieces []Piece
}

// --X-
// --X-
// --X-
// --X-
var pieceLongBoi = Piece{Point{0, -2}, Point{0, -1}, Point{0, 0}, Point{0, 1}}

/*
// ----  ----
// -X--  --X-
// -X--  --X-
// -XX-  -XX-
var pieceL = Piece{Point{-1, -1}, Point{-1, 1}, Point{-1, 2}, Point{1, 2}}
var pieceInvertedL = Piece{Point{1, -1}, Point{1, 1}, Point{-1, 2}, Point{1, 2}}

// ----  ----
// -X--  --X-
// -XX-  -XX-
// --X-  -X--
var pieceZ = Piece{Point{-1,-1}, Point{-1,1}, Point{1,1}, Point{1,2}}
var pieceInvertedZ = Piece{Point{1,-1}, Point{-1,1}, Point{1,1}, Point{-1,2}}

// ----
// -XX-
// -XX-
// ----
var pieceBlock = Piece{Point{-1,-1}, Point{-1,1}, Point{1,-1}, Point{1,1}}

// ----
// -X--
// XXX-
// ----
var pieceTBoi = Piece{Point{-1,-1}, Point{-2,1}, Point{-1,1}, Point{1,1}}
*/

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
