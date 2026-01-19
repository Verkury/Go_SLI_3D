package engine

import (
	. "math"
)

type Point3D struct {
	x, y, z float64
}

const symbols string = "•○●#@"

func (s Screen) AddPoint(x, y int) {
    if y >= 0 && y < len(s) && x >= 0 && x < len(s[0]) {
        s[y][x] = "@" 
    }
}

// I brazenly stole it
// But I tried to do this
func (s Screen) AddLine(startX, startY, endX, endY int) {
	dx := Abs(float64(endX - startX))
	dy := Abs(float64(endY - startY))
	sx := sign(endX-startX)
	sy := sign(endY-startY)
	err := dx - dy

	x, y := startX, startY

	for  {
		s.AddPoint(x,y)
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

func (s Screen) AddRect(x0,y0,x1,y1 int) {
	s.AddLine(x0, y0, x1, y0)
	s.AddLine(x0, y0, x0, y1)
	s.AddLine(x1, y0, x1, y1)
	s.AddLine(x0, y1, x1, y1)
}

func (s Screen) AddSquare(x0, y0, size int, filed bool) {
	x1, y1 := x0+size*Aspect, y0+size

	s.AddLine(x0, y0, x1, y0)
	s.AddLine(x0, y0, x0, y1)
	s.AddLine(x1, y0, x1, y1)
	s.AddLine(x0, y1, x1, y1)

	if filed {
		for line := y0; line < y1; line++ {
			s.AddLine(x0, line, x1, line)
		}
	}
}

func (s Screen) AddCube(x, y, z, size int, angleX, angleY, angleZ float64) {
    x1, y1, z1 := x + size, y + size, z + size

    vertices := []Point3D{
        {float64(x),  float64(y),  float64(z)},  // 0
        {float64(x1), float64(y),  float64(z)},  // 1
        {float64(x1), float64(y1), float64(z)},  // 2
        {float64(x),  float64(y1), float64(z)},  // 3
        {float64(x),  float64(y),  float64(z1)}, // 4
        {float64(x1), float64(y),  float64(z1)}, // 5
        {float64(x1), float64(y1), float64(z1)}, // 6
        {float64(x),  float64(y1), float64(z1)}, // 7
    }

    edges := [][2]int{
        {0, 1}, {1, 2}, {2, 3}, {3, 0}, 
        {4, 5}, {5, 6}, {6, 7}, {7, 4}, 
        {0, 4}, {1, 5}, {2, 6}, {3, 7}, 
    }
    
    w, h := len((s)[0]), len(s)
	for _, edge := range edges {
		p1 := Rotate(vertices[edge[0]], angleX, angleY, angleZ)
		p2 := Rotate(vertices[edge[1]], angleX, angleY, angleZ)

		x1, y1 := pointPerspectiveProjection(p1, w, h)
		x2, y2 := pointPerspectiveProjection(p2, w, h)

		s.AddLine(x1, y1, x2, y2)
	}
}

// I stole that too.
func Rotate(p Point3D, angleX, angleY, angleZ float64) Point3D {
	rad := angleX
	p = Point3D{p.x, p.y*Cos(rad) - p.z*Sin(rad), p.y*Sin(rad) + p.z*Cos(rad)}
	rad = angleY
	p = Point3D{p.x*Cos(rad) + p.z*Sin(rad), p.y, -p.x*Sin(rad) + p.z*Cos(rad)}
	rad = angleZ
	p = Point3D{p.x*Cos(rad) - p.y*Sin(rad), p.x*Sin(rad) + p.y*Cos(rad), p.z}
	return p
}