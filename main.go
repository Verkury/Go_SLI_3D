package main

import (
	"math/rand/v2"

	"github.com/Verkury/Go_SLI_3D/engine"
)

func main() {
	scene := engine.Init()
	go scene.Start()

	angleX, angleY, angleZ := 0.0, 0.0, 0.0

	for {
		scene.Scr.ClearScrean()

		var LightSource = engine.Point3D{X: -1, Y: -1, Z: -1}

		scene.Scr.AddCube(-1, -1, -1, 2, angleX, angleY, angleZ, LightSource)

		min := 1
		max := 6
		angleX += float64(rand.IntN(max-min) + min)/100
		angleY += float64(rand.IntN(max-min) + min)/100
		angleZ += float64(rand.IntN(max-min) + min)/100
		scene.Draw()
	}
}