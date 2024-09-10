package main

import (
	"sync"

	"github.com/gdamore/tcell/v2"
)

const GameFrameWidth = 80
const GameFrameHigh = 25
const GameFrameSymbol = '║'
const BulletSymbol = '*'

var screen tcell.Screen

var player *GameObject
var zombies []*GameObject
var bullets []*GameObject
var pointsToClear []*Point
var isGamePaused bool
var isGameOver bool
var debugLog string
var score int
var inputs []string
var mu sync.Mutex
var restart bool
