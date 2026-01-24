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

type Screen struct {
	Chars [][]string
	ZBuffer [][]float64
}

type Scene struct {
	Scr Screen
	frames int
	lastfps int64
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

	chars := make([][]string, height)
	zbuf := make([][]float64, height)
    for i := range chars {
        chars[i] = make([]string, width)
		zbuf[i] = make([]float64, width)
    }

	s := Screen{Chars: chars, ZBuffer: zbuf}
	s.ClearScrean()
	
	return Scene{s, 0, 0}
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
	winX := int(float64(width)/2 + p.X*Scale*Aspect)
	winY := int(float64(height)/2 + p.Y*Scale)
	return winX, winY
}