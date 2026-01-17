package engine

import (
	. "math"
)

type Point3D struct {
	x, y, z float64
}

const symbols string = "•○●#@"

func (s Screen) addPoint(x, y int) {
    if y >= 0 && y < len(s) && x >= 0 && x < len(s[0]) {
        s[y][x] = "@" 
    }
}

// I brazenly stole it
// But I tried to do this
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

func (s Screen) addSquare(x0, y0, size int, filed bool) {
	x1, y1 := x0+size*Aspect, y0+size

	s.addLine(x0, y0, x1, y0)
	s.addLine(x0, y0, x0, y1)
	s.addLine(x1, y0, x1, y1)
	s.addLine(x0, y1, x1, y1)

	if filed {
		for line := y0; line < y1; line++ {
			s.addLine(x0, line, x1, line)
		}
	}
}

// I stole that too.
func rotate(p Point3D, angleX, angleY, angleZ float64) Point3D {
	rad := angleX
	p = Point3D{p.x, p.y*Cos(rad) - p.z*Sin(rad), p.y*Sin(rad) + p.z*Cos(rad)}
	rad = angleY
	p = Point3D{p.x*Cos(rad) + p.z*Sin(rad), p.y, -p.x*Sin(rad) + p.z*Cos(rad)}
	rad = angleZ
	p = Point3D{p.x*Cos(rad) - p.y*Sin(rad), p.x*Sin(rad) + p.y*Cos(rad), p.z}
	return p
}