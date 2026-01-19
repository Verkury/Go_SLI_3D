package engine

import (
	"os"
	"sync"

	"golang.org/x/term"
)

var wg sync.WaitGroup

var FPS = 60
var DEALAY = 50
var SHOWFPS = true

var Scale = 8.0

const Aspect = 2.0

type Screen [][]string

type Scene struct {
	Scr Screen
	frames int
	lastfps int
}

func getTerminalSize() (int, int){
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 14
	}
	return width, height
}

func Init() Scene {
	width, height := getTerminalSize()

	screen := make(Screen, height)
    for i := range screen {
        screen[i] = make([]string, width)
    }

	screen.ClearScrean()
	
	return Scene{screen, 0, 0}
}

func (s *Scene) Start() {
	wg.Add(1)

	go updateStream(s)

	defer wg.Wait()
	
}

func sign(n int) int {
	if n > 0 { return 1 }
	if n < 0 { return -1 }
	return 0
}

func pointPerspectiveProjection(p Point3D, width, height int) (int, int) {
	winX := int(float64(width)/2 + p.x*Scale*Aspect)
	winY := int(float64(height)/2 + p.y*Scale)
	return winX, winY
}