package main

type Point struct {
	row, col int
	symbol   rune
}

type GameObject struct {
	points         []*Point
	velRow, velCol int
}

type Color int

const (
	White Color = iota
	Black
	Blue
	Red
	Green
	Yellow
)
