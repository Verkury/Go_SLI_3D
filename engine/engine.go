package engine

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"math/rand/v2"

	"golang.org/x/term"
)

var wg sync.WaitGroup

var FPS = 60
var DEALAY = 50
var SHOWFPS = true

var Scale = 8.0

const Aspect = 2.0

type Screen [][]string

func getTerminalSize() (int, int){
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 14
	}
	return width, height
}

func Start() {
	width, height := getTerminalSize()

	screan := make(Screen, height)
    for i := range screan {
        screan[i] = make([]string, width)
    }

	clearScrean := makeClearScrean(screan)
	screan = clearScrean
	/*
	screan.addLine(3 * Aspect,3,13 * Aspect,13)
	screan.addRect(30 * Aspect,30, 40 * Aspect,40)
	screan.addSquare(5 * Aspect,5,4, true)
	*/
	wg.Add(1)
	go updateStream(&screan)

	defer wg.Wait()
	
}

func (s Screen) drow() {
	var out strings.Builder
	
	out.WriteString("\033[H") 

	for y := range s {
		out.WriteString("\n")
		for x := range s[y] {
			if s[y][x] == "" {
				out.WriteString(" ")
			} else {
				out.WriteString(s[y][x])
			}
		}
	}

	fmt.Print(out.String())
	
}

func updateStream(s *Screen) {
	defer wg.Done()
	
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	
	frames := 0
	fpsDisplay := 0
	angleX, angleY, angleZ := 0.0, 0.0, 0.0

	vertices := []Point3D{
        {-1, -1, -1}, {1, -1, -1}, {1, 1, -1}, {-1, 1, -1},
        {-1, -1, 1}, {1, -1, 1}, {1, 1, 1}, {-1, 1, 1},
    }
    
    edges := [][2]int{
        {0, 1}, {1, 2}, {2, 3}, {3, 0},
        {4, 5}, {5, 6}, {6, 7}, {7, 4},
        {0, 4}, {1, 5}, {2, 6}, {3, 7},
    }

	for {
		select {
		case <-ticker.C:
			fpsDisplay = frames
			frames = 0
		default:
			frames++
			min := 1
			max := 6
			angleX += float64(rand.IntN(max-min) + min)/100
			angleY += float64(rand.IntN(max-min) + min)/100
			angleZ += float64(rand.IntN(max-min) + min)/100
			makeClearScrean(*s)

			w, h := len((*s)[0]), len(*s)
			for _, edge := range edges {
				p1 := rotate(vertices[edge[0]], angleX, angleY, angleZ)
				p2 := rotate(vertices[edge[1]], angleX, angleY, angleZ)

				x1, y1 := pointPerspectiveProjection(p1, w, h)
				x2, y2 := pointPerspectiveProjection(p2, w, h)

				s.addLine(x1, y1, x2, y2)
			}
			
			if SHOWFPS {
				s.renderFPS(fpsDisplay*4)
			}
			
			s.drow()
			
			time.Sleep(time.Duration(1000/FPS) * time.Millisecond)
		}
	}
}

func (s Screen) renderFPS(fps int) {
	str := fmt.Sprintf("FPS: %d", fps)
	for i, char := range str {
		if i < len(s[0]) {
			s[0][i] = string(char)
		}
	}
}

func makeClearScrean(screan Screen) Screen{
	for y := range screan {
		for x := range screan[y] {
			screan[y][x] = " "
		}
	}
	return screan
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