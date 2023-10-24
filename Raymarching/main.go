package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 500
	screenHeight = 500
	screenDepth  = 100
)

//go:embed fragment.go
var shaderBytes []byte

type Game struct {
	fragmentShader ebiten.Shader
	time           int

	cameraPosition Vector3
	iterationMax   int
}

func (g *Game) Update() error {
	g.time++
	return nil
}

func (g *Game) DrawShsader(screen *ebiten.Image) {
	op := &ebiten.DrawRectShaderOptions{}
	cx, cy := ebiten.CursorPosition()

	op.Uniforms = map[string]any{
		"Time":   float64(g.time),
		"Cursor": []float32{float32(cx), float32(cy)},
	}
	op.Images[0] = ebiten.NewImage(screenWidth, screenHeight)

	screen.DrawRectShader(screenWidth, screenHeight, &g.fragmentShader, op)
}

func (g *Game) GetMinDistance(position Vector3) float64 {
	spherePosition := Vector3{250, 250, 100}
	sphereRadius := 50

	distanceToSphere := math.Sqrt(math.Pow(spherePosition.x-position.x, 2)+math.Pow(spherePosition.y-position.y, 2)+math.Pow(spherePosition.z-position.z, 2)) - float64(sphereRadius)

	return distanceToSphere
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.Fill(color.Black)
	//g.DrawShsader(screen)

	//screen.Set(100, 100, color.White)
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			screenPixelPosition := Vector3{float64(x), float64(y), screenDepth}
			directionX := screenPixelPosition.x - g.cameraPosition.x
			directionY := screenPixelPosition.y - g.cameraPosition.y
			directionZ := screenPixelPosition.z - g.cameraPosition.z

			directionVectorLength := math.Sqrt(directionX*directionX + directionY*directionY + directionZ*directionZ)
			directionVector := Vector3{directionX / directionVectorLength, directionY / directionVectorLength, directionZ / directionVectorLength}

			currentPosition := Vector3{g.cameraPosition.x, g.cameraPosition.y, g.cameraPosition.z}

			for iteration := 0; iteration < g.iterationMax; iteration++ {
				minDist := g.GetMinDistance(currentPosition)
				//fmt.Println(minDist)

				if minDist < 0.0001 {
					// collision
					screen.Set(x, y, color.White)
				}

				if minDist > 1000.0 {
					// too far nothing here
					// miss
					break
				}

				// travel along vector
				currentPosition.x += directionVector.x * minDist
				currentPosition.y += directionVector.y * minDist
				currentPosition.z += directionVector.z * minDist
			}

		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Noise (Ebitengine Demo)")

	g := &Game{}
	g.cameraPosition = Vector3{250, 250, 0}
	g.iterationMax = 5

	s, err := ebiten.NewShader(shaderBytes)
	if err != nil {
		log.Fatal(err)
	}
	g.fragmentShader = *s
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
