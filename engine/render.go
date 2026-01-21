package engine

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

func (s *Scene) Draw() {
	var out strings.Builder

	if SHOWFPS {
		s.Scr.renderFPS(int(atomic.LoadInt64(&s.lastfps)))
	}

	out.WriteString("\033[H") 

	for y := range s.Scr {
		out.WriteString("\n")
		for x := range s.Scr[y] {
			if s.Scr[y][x] == "" {
				out.WriteString(" ")
			} else {
				out.WriteString(s.Scr[y][x])
			}
		}
	}

	fmt.Print(out.String())
	s.frames++ 
	time.Sleep(time.Duration(1000/FPS) * time.Millisecond)
}


func updateStream(s *Scene) {
	defer wg.Done()
	
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			atomic.StoreInt64(&s.lastfps, int64(s.frames))
			s.frames = 0
		default:		
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

func (screan Screen) ClearScrean() {
	for y := range screan {
		for x := range screan[y] {
			screan[y][x] = " "
		}
	}
}