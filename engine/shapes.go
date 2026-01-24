package engine

import (
	. "math"
)

type Point3D struct {
	X, Y, Z float64
}

const symbols string = ".:-=+*#%@"

func (s Screen) AddPoint(x, y int, z float64, char string) {
    if y >= 0 && y < len(s.Chars) && x >= 0 && x < len(s.Chars[0]) {
        if z < s.ZBuffer[y][x] {
			s.Chars[y][x] = char
			s.ZBuffer[y][x] = z
		}
    }
}

func (s Screen) AddCube(x, y, z, size int, angleX, angleY, angleZ float64, lightSource Point3D) {
	x1, y1, z1 := float64(x), float64(y), float64(z)
	s2 := float64(size)

	v := []Point3D{
		{x1, y1, z1}, {x1 + s2, y1, z1}, {x1 + s2, y1 + s2, z1}, {x1, y1 + s2, z1},
		{x1, y1, z1 + s2}, {x1 + s2, y1, z1 + s2}, {x1 + s2, y1 + s2, z1 + s2}, {x1, y1 + s2, z1 + s2},
	}

	faces := [][]int{
		{0, 1, 2, 3}, {7, 6, 5, 4}, // перед/зад
		{0, 4, 5, 1}, {2, 6, 7, 3}, // верх/низ
		{0, 3, 7, 4}, {1, 5, 6, 2}, // бока
	}

	w, h := len(s.Chars[0]), len(s.Chars)
	light := lightSource.Normalize()

	for _, f := range faces {
		p0 := Rotate(v[f[0]], angleX, angleY, angleZ)
		p1 := Rotate(v[f[1]], angleX, angleY, angleZ)
		p2 := Rotate(v[f[2]], angleX, angleY, angleZ)
		p3 := Rotate(v[f[3]], angleX, angleY, angleZ)

		normal := CalculateNormal(p0, p1, p2).Normalize()

		if normal.Z > 0 { continue }

		shade := getShade(normal, light)

		s.FillTriangle(p0, p1, p2, shade, w, h)
		s.FillTriangle(p0, p2, p3, shade, w, h)
		
	}
}

// I stole that too.
func Rotate(p Point3D, angleX, angleY, angleZ float64) Point3D {
	rad := angleX
	p = Point3D{p.X, p.Y*Cos(rad) - p.Z*Sin(rad), p.Y*Sin(rad) + p.Z*Cos(rad)}
	rad = angleY
	p = Point3D{p.X*Cos(rad) + p.Z*Sin(rad), p.Y, -p.X*Sin(rad) + p.Z*Cos(rad)}
	rad = angleZ
	p = Point3D{p.X*Cos(rad) - p.Y*Sin(rad), p.X*Sin(rad) + p.Y*Cos(rad), p.Z}
	return p
}

func CalculateNormal(p1, p2, p3 Point3D) Point3D {
	v1 := Point3D{p2.X - p1.X, p2.Y - p1.Y, p2.Z - p1.Z}
	v2 := Point3D{p3.X - p1.X, p3.Y - p1.Y, p3.Z - p1.Z}

	return Point3D{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: v1.Z*v2.X - v1.X*v2.Z,
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

func (p Point3D) Normalize() Point3D {
	mag := Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
	if mag == 0 { return Point3D{0,0,0} }
	return Point3D{p.X / mag, p.Y / mag, p.Z / mag}
}

func getShade(normal, light Point3D) string {
	dot := normal.X*light.X + normal.Y*light.Y + normal.Z*light.Z
	
	if dot < 0.1 { dot = 0.1 }
    
	index := int(dot * float64(len(symbols)-1))

	if index < 0 { index = 0 }
    if index >= len(symbols) { index = len(symbols) - 1 }

	return string(symbols[index])
}

func (s Screen) DrowShadedLine(p1, p2 Point3D, shade string, z float64, width, height int) {
	startX, startY := pointPerspectiveProjection(p1, width, height)
	endX, endY := pointPerspectiveProjection(p2, width, height)

	dx := Abs(float64(endX - startX))
	dy := Abs(float64(endY - startY))
	sx := sign(endX - startX)
	sy := sign(endY - startY)
	err := dx - dy

	x, y := startX, startY

	for {
		s.AddPoint(x, y, z, shade)
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


func (s Screen) FillTriangle(p1, p2, p3 Point3D, shade string, width, height int) {
	x1, y1 := pointPerspectiveProjection(p1, width, height)
	x2, y2 := pointPerspectiveProjection(p2, width, height)
	x3, y3 := pointPerspectiveProjection(p3, width, height)

	avgZ := (p1.Z + p2.Z + p3.Z) / 3

	if y1 > y2 { x1, x2 = x2, x1; y1, y2 = y2, y1 }
	if y1 > y3 { x1, x3 = x3, x1; y1, y3 = y3, y1 }
	if y2 > y3 { x2, x3 = x3, x2; y2, y3 = y3, y2 }

	interpolate := func(y, yStart, yEnd, xStart, xEnd int) int {
		if yStart == yEnd { return xStart }
		return xStart + (xEnd-xStart)*(y-yStart)/(yEnd-yStart)
	}

	for y := y1; y <= y3; y++ {
		var startX, endX int
		if y < y2 {
			startX = interpolate(y, y1, y2, x1, x2)
			endX = interpolate(y, y1, y3, x1, x3)
		} else {
			startX = interpolate(y, y2, y3, x2, x3)
			endX = interpolate(y, y1, y3, x1, x3)
		}

		if startX > endX { startX, endX = endX, startX }

		for x := startX; x <= endX; x++ {
			s.AddPoint(x, y, avgZ, shade)
		}
	}
}