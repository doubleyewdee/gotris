package main

import (
	"github.com/gdamore/tcell"
)

const (
	BOARD_WIDTH  = 10
	BOARD_HEIGHT = 20
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