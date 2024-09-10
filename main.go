package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

func main() {
	for {
		InitScreen()
		InitGameState()
		inputchan := InitUserInput()
		ReadInput(inputchan)

		for !isGameOver {
			//ClearScreen()
			ProcessInputs()
			UpdateState()
			DrawState()

			time.Sleep(75 * time.Millisecond)
		}

		DrawGameOver()

		// Wait for the user input after the game is over
		for isGameOver && !restart {
			ProcessInputs()
			time.Sleep(75 * time.Millisecond)
		}

		// Clean resources
		screen.Fini()

		if !restart {
			break
		}
	}
	// Clean resources
	screen.Fini()
}

func UpdateState() {
	if isGamePaused {
		return
	}

	UpdateZombies()
	UpdateBullets()
	ColisionDetection()
}

func UpdateZombies() {
	MoveGameObjects(zombies)

	chance := rand.Intn(100)
	if chance < 5 {
		SpawnZombie()
	}
}

func UpdateBullets() {
	MoveGameObjects(bullets)
}

func SpawnZombie() {
	originRow, originCol := rand.Intn(GameFrameHigh-3), GameFrameWidth-3
	zombies = append(zombies, &GameObject{
		points: []*Point{
			{row: originRow, col: originCol, symbol: '0'},
			{row: originRow + 1, col: originCol, symbol: '|'},
			{row: originRow + 1, col: originCol - 1, symbol: '\\'},
			{row: originRow + 2, col: originCol, symbol: '|'},
			{row: originRow + 3, col: originCol - 1, symbol: '/'},
			{row: originRow + 3, col: originCol + 1, symbol: '\\'},
		},
		velRow: 0, velCol: -1,
	})
}

func SpawnBullet() {
	rowOrigin, colOrigin := player.points[5].row, player.points[5].col
	bullets = append(bullets, &GameObject{
		points: []*Point{
			{row: rowOrigin, col: colOrigin, symbol: BulletSymbol},
		},
		velRow: 0, velCol: 2,
	})
}

func MoveGameObjects(objs []*GameObject) {
	for _, obj := range objs {
		for _, p := range obj.points {
			copy := *p
			pointsToClear = append(pointsToClear, &copy)
			p.col += obj.velCol
			p.row += obj.velRow
		}
	}
}

func ColisionDetection() {
	// Zombies with wall - Zoombies with player
	for _, z := range zombies {
		// Colided with wall
		for _, zp := range z.points {
			if zp.col <= 0 {
				isGameOver = true
			}
		}

		// Colided with player
		if AreObjectsCollinding(z, player, 1) {
			isGameOver = true
		}
	}

	// Bullets with zombies - Bullets with wall
	for bi := len(bullets) - 1; bi >= 0; bi-- {
		for zi := len(zombies) - 1; zi >= 0; zi-- {
			if AreObjectsCollinding(bullets[bi], zombies[zi], 1) {
				bullets = append(bullets[:bi], bullets[bi+1:]...)
				zombies = append(zombies[:zi], zombies[zi+1:]...)
				score++
				break
			}
		}

		if bi < len(bullets) && bullets[bi].points[0].col >= GameFrameWidth {
			bullets = append(bullets[:bi], bullets[bi+1:]...)
		}
	}
}

func AreObjectsCollinding(obj1, obj2 *GameObject, radius int) bool {
	for _, p1 := range obj1.points {
		for _, p2 := range obj2.points {
			if math.Abs(float64(p1.col-p2.col)) <= float64(radius) && p1.row == p2.row {
				return true
			}
		}
	}

	return false
}

func ClearScreen() {

	for _, p := range pointsToClear {
		DrawInsideGameFrame(p.row, p.col, 1, 1, ' ')
	}
	pointsToClear = []*Point{}

	// for _, z := range zombies {
	// 	for _, p := range z.points {
	// 		DrawInsideGameFrame(p.row, p.col, 1, 1, ' ')
	// 	}
	// }

	// for _, b := range bullets {
	// 	for _, p := range b.points {
	// 		DrawInsideGameFrame(p.row, p.col, 1, 1, ' ')
	// 	}
	// }

	// for _, p := range player.points {
	// 	DrawInsideGameFrame(p.row, p.col, 1, 1, ' ')
	// }
}

func DrawScore() {
	row, col := GetGameFrameTopLeft()
	PrintString(row-2, col, fmt.Sprintf("Score: %d", score))
}

func DrawState() {
	if isGamePaused {
		return
	}

	//screen.Clear() // TODO - Improve it to clean only cells that needed to be cleanned
	ClearScreen()
	DrawScore()
	DrawGameFrame()
	DrawGameObjects(append(append([]*GameObject{player}, zombies...), bullets...))

	PrintString(0, 0, debugLog)
	screen.Show()
}

func DrawGameObjects(objs []*GameObject) {
	for _, obj := range objs {
		for _, p := range obj.points {
			DrawInsideGameFrame(p.row, p.col, 1, 1, p.symbol)
		}
	}
}

func DrawInsideGameFrame(row, col, width, height int, ch rune, color ...Color) {
	rowOffset, colOffset := GetGameFrameTopLeft()
	DrawFilledRect(row+rowOffset, col+colOffset, width, height, ch, color...)
}

func DrawFilledRect(row, col, width, height int, ch rune, color ...Color) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			var style tcell.Style
			if color == nil {
				style = tcell.StyleDefault
			} else {
				style = GetColor(color[0])
			}
			screen.SetContent(col+c, row+r, ch, nil, style)
		}
	}
}

func DrawUnfilledRect(row, col, width, height int, ch rune, color ...Color) {
	var style tcell.Style
	if color == nil {
		style = tcell.StyleDefault
	} else {
		style = GetColor(color[0])
	}

	for c := 0; c < width; c++ {
		screen.SetContent(col+c, row, ch, nil, style)
		screen.SetContent(col+c, row+height-1, ch, nil, style)
	}

	for r := 0; r < height-1; r++ {
		screen.SetContent(col, row+r, ch, nil, style)
		screen.SetContent(col+width-1, row+r, ch, nil, style)
	}
}

func DrawGameFrame() {
	gameFrameTopLeftRow, gameFrameTopLeftCol := GetGameFrameTopLeft()
	row, col := gameFrameTopLeftRow-1, gameFrameTopLeftCol-1
	width, height := GameFrameWidth+2, GameFrameHigh+2

	DrawUnfilledRect(row, col, width, height, GameFrameSymbol)
	//DrawUnfilledRect(row+1, col+1, GameFrameWidth, GameFrameHigh, '*')
}

func DrawGameOver() {
	screenWidth, screenHeight := screen.Size()
	PrintStringCentered(screenHeight/2, screenWidth/2, "Game Over!")
	PrintStringCentered(screenHeight/2+1, screenWidth/2+1, fmt.Sprint("Your score is: ", score))
	PrintStringCentered(screenHeight/2+2, screenWidth/2+2, "Press Esc or 'q' to leave")
	PrintStringCentered(screenHeight/2+3, screenWidth/2+3, "Press Enter to try again!")
	screen.Show()
}

func PrintString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func PrintStringCentered(row, col int, str string) {
	col = col - len(str)/2
	PrintString(row, col, str)
}

func GetGameFrameTopLeft() (int, int) {
	screnWidth, screenHeight := screen.Size()
	return (screenHeight - GameFrameHigh) / 2, (screnWidth - GameFrameWidth) / 2
}

func InitScreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	screen.HideCursor()
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorNone)
	screen.SetStyle(defStyle)
}

func InitGameState() {
	player = &GameObject{
		points: []*Point{
			{row: 5, col: 1, symbol: '0'},
			{row: 6, col: 1, symbol: '|'},
			{row: 7, col: 1, symbol: '|'},
			{row: 6, col: 2, symbol: '-'},
			{row: 6, col: 3, symbol: '-'},
			{row: 6, col: 4, symbol: '-'},
			{row: 8, col: 0, symbol: '/'},
			{row: 8, col: 2, symbol: '\\'},
		},
	}
	zombies = nil
	score = 0
	restart = false
	isGameOver = false
	inputs = make([]string, 0, 100)
}

func InitUserInput() chan string {
	inputChan := make(chan string) // Non-buffered channel
	go func() {
		for {

			switch ev := screen.PollEvent().(type) { // Block waiting for the event
			case *tcell.EventKey:
				debugLog = ev.Name()
				inputChan <- ev.Name()
			}
		}
	}()
	return inputChan
}

func ReadInput(inputChan chan string) {
	go func() {
		for {
			key := <-inputChan // Wait until has something to read. It locks while waiting because it is a non-buffered channel
			mu.Lock()
			inputs = append(inputs, key)
			mu.Unlock()
		}
	}()
}

func ProcessInputs() {
	mu.Lock()
	// I still not sure which one is better... process all inputs and the render it or process only the last input and render it.
	for i := len(inputs) - 1; i >= 0; i-- {
		HandleUserInput(inputs[i])
	}
	// if len(inputs) > 0 {
	// 	HandleUserInput(inputs[len(inputs)-1])
	// }

	inputs = inputs[:0] // Clean the buffer of inputs
	mu.Unlock()
}

func HandleUserInput(key string) {
	if key == "Rune[q]" || key == "Esc" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Enter" && isGameOver {
		restart = true
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	} else if (key == "Rune[w]" || key == "Up") && !IsObjectHittingWall(player, -1, 0) {
		MovePlayer(-1, 0)
	} else if (key == "Rune[s]" || key == "Down") && !IsObjectHittingWall(player, 1, 0) {
		MovePlayer(1, 0)
	} else if (key == "Rune[a]" || key == "Left") && !IsObjectHittingWall(player, 0, -1) {
		MovePlayer(0, -1)
	} else if (key == "Rune[d]" || key == "Right") && !IsObjectHittingWall(player, 0, 1) {
		MovePlayer(0, 1)
	} else if key == "Rune[ ]" || key == "Enter" {
		debugLog = "Shoot"
		SpawnBullet()
	}
}

func MovePlayer(velRow, velCol int) {
	for _, p := range player.points {
		copy := *p
		pointsToClear = append(pointsToClear, &copy)
		p.col += velCol
		p.row += velRow
	}
}

func IsObjectHittingWall(obj *GameObject, velRow, velCol int) bool {
	for _, p := range obj.points {
		if p.col+velCol < 0 || p.col+velCol >= GameFrameWidth || p.row+velRow < 0 || p.row+velRow >= GameFrameHigh {
			return true
		}
	}

	return false
}
