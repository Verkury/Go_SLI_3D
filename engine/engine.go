package engine

import (
	"fmt"
	. "math"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

var wg sync.WaitGroup

var FPS = 60
var SHOWFPS = true

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

	screan.addLine(3,3,13,13)
	screan.addRect(30,30, 40,40)

	wg.Add(1)
	go updateStream(&screan)

	defer wg.Wait()
	
}

func (s Screen) addPoint(x, y int) {
    renderX := int(float64(x) * Aspect)

    if y >= 0 && y < len(s) && renderX >= 0 && renderX < len(s[0]) {
        s[y][renderX] = "@" 
    }
}

func (s Screen) addLine(startX, startY, endX, endY int) {
	dx := Abs(float64(endX - startX))
	dy := Abs(float64(endY - startY))
	sx := sign(endX-startX)
	sy := sign(endY-startY)
	err := dx - dy

	x, y := startX, startY

	for  {
		s.addPoint(x,y)
		if x == endX && y == endY {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func (s Screen) addRect(x0,y0,x1,y1 int) {
	s.addLine(x0, y0, x1, y0)
	s.addLine(x0, y0, x0, y1)
	s.addLine(x1, y0, x1, y1)
	s.addLine(x0, y1, x1, y1)
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

	for {
		select {
		case <-ticker.C:
			fpsDisplay = frames
			frames = 0
		default:
			frames++
			
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
			s[1][i] = string(char)
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