package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	_screenWidth  = 800
	_screenHeight = 800
)

var (
	whiteImage = ebiten.NewImage(3, 3)
)

func init() {
	whiteImage.Fill(color.White)
}

type Game struct {
	positionX []float64
	positionY []float64

	directionX []float64
	directionY []float64
}

func (g *Game) NormalizeDirections() {
	for i := range g.directionX {
		length := math.Sqrt(g.directionX[i]*g.directionX[i] + g.directionY[i]*g.directionY[i])

		g.directionX[i] /= length
		g.directionY[i] /= length
	}
}

const (
	maxForce     = 0.0001 // Maximum steering force
	desiredSep   = 100.0  // Desired separation between boids
	neighborDist = 200.0  // Distance to consider other boids as neighbors
)

func (g *Game) Update() error {
	g.NormalizeDirections()

	for i := range g.positionX {
		// Initialize steering vectors to zero
		steerSeparation := [2]float64{0, 0}
		steerAlignment := [2]float64{0, 0}
		steerCohesion := [2]float64{0, 0}

		countAlignment := 0
		countCohesion := 0

		for j := range g.positionX {
			if i == j {
				continue
			}

			dx := g.positionX[j] - g.positionX[i]
			dy := g.positionY[j] - g.positionY[i]
			distance := math.Sqrt(dx*dx + dy*dy)

			// Separation
			if distance < desiredSep {
				diffX := g.positionX[i] - g.positionX[j]
				diffY := g.positionY[i] - g.positionY[j]
				normalizeFactor := 1.0 / (distance + 0.001) // Add small value to avoid division by zero
				steerSeparation[0] += diffX * normalizeFactor
				steerSeparation[1] += diffY * normalizeFactor
			}

			// Alignment and Cohesion
			if distance < neighborDist {
				steerAlignment[0] += g.directionX[j]
				steerAlignment[1] += g.directionY[j]
				countAlignment++

				steerCohesion[0] += g.directionX[j]
				steerCohesion[1] += g.directionY[j]
				countCohesion++
			}
		}

		// Average alignment and cohesion
		if countAlignment > 0 {
			steerAlignment[0] /= float64(countAlignment)
			steerAlignment[1] /= float64(countAlignment)
		}

		if countCohesion > 0 {
			steerCohesion[0] /= float64(countCohesion)
			steerCohesion[1] /= float64(countCohesion)
			steerCohesion[0] = steerCohesion[0] - g.directionX[i]
			steerCohesion[1] = steerCohesion[1] - g.directionY[i]
		}

		// Combine the rules (you might want to weight these differently)
		g.directionX[i] += (steerSeparation[0] + steerAlignment[0] + steerCohesion[0]) * 0.1
		g.directionY[i] += (steerSeparation[1] + steerAlignment[1] + steerCohesion[1]) * 0.1

		// Limit the magnitude of direction to maxForce
		length := math.Sqrt(g.directionX[i]*g.directionX[i] + g.directionY[i]*g.directionY[i])
		if length > maxForce {
			g.directionX[i] /= length
			g.directionY[i] /= length
			g.directionX[i] *= maxForce
			g.directionY[i] *= maxForce
		}
	}

	g.NormalizeDirections()

	velocity := 1.0
	for i := range g.positionX {
		g.positionX[i] += g.directionX[i] * velocity
		g.positionY[i] += g.directionY[i] * velocity

		if g.positionX[i] <= 0 || g.positionX[i] >= _screenWidth {
			g.positionX[i] = _screenWidth / 2
			g.directionX[i] *= -1
		}

		if g.positionY[i] <= 0 || g.positionY[i] >= _screenHeight {
			g.positionY[i] = _screenHeight / 2
			g.directionY[i] *= -1
		}
	}

	return nil
}

func GenerateVertices(x, y, directionX, directionY float64) []ebiten.Vertex {

	// Calculate the rotation angle
	theta := math.Atan2(directionY, directionX) + math.Pi/2

	vs := []ebiten.Vertex{}
	size := 20

	// Helper function to rotate a point around the origin
	rotate := func(px, py, theta float64) (float64, float64) {
		return px*math.Cos(theta) - py*math.Sin(theta), px*math.Sin(theta) + py*math.Cos(theta)
	}

	// List of points relative to the center of the shape
	points := [][]float64{
		{0, -float64(size) / 2},
		{float64(size) / 2, float64(size) / 2},
		{0, 0},
		{-float64(size) / 2, float64(size) / 2},
	}

	for _, p := range points {
		// Rotate each point
		rx, ry := rotate(p[0], p[1], theta)

		vs = append(vs, ebiten.Vertex{
			DstX:   float32(x + rx),
			DstY:   float32(y + ry),
			SrcX:   0,
			SrcY:   0,
			ColorR: float32(1),
			ColorG: float32(1),
			ColorB: float32(1),
			ColorA: 1,
		})
	}

	return vs
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawTrianglesOptions{AntiAlias: false}
	indices := []uint16{0, 1, 2, 2, 3, 0}
	lineLength := 50.0

	for i := range g.positionX {
		//log.Printf("%f %f", g.positionX[i], g.positionY[i])
		screen.DrawTriangles(GenerateVertices(g.positionX[i], g.positionY[i], g.directionX[i], g.directionY[i]), indices, whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image), op)
		// Assuming you are in the loop where you draw your polygons

		x1 := g.positionX[i] + g.directionX[i]*lineLength
		y1 := g.positionY[i] + g.directionY[i]*lineLength

		vector.StrokeLine(screen, float32(g.positionX[i]), float32(g.positionY[i]), float32(x1), float32(y1), 1, color.RGBA{0, 255, 0, 255}, false)

		//vector.DrawFilledRect(screen, float32(g.positionX[i]), float32(g.positionY[i]), 1, 1, color.White, false)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return _screenWidth, _screenHeight
}

func main() {
	ebiten.SetWindowSize(_screenWidth, _screenHeight)
	ebiten.SetWindowTitle("Hello, World!")

	game := Game{}

	for i := 0; i < 100; i++ {
		maxOffset := 160.0
		game.positionX = append(game.positionX, maxOffset+rand.Float64()*(_screenWidth-2*maxOffset))
		game.positionY = append(game.positionY, maxOffset+rand.Float64()*(_screenHeight-2*maxOffset))

		game.directionX = append(game.directionX, 2*rand.Float64()-1)
		game.directionY = append(game.directionY, 2*rand.Float64()-1)
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
