package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 1000
	screenHeight = 1000

	deltaT             = 0.02
	attractionDistance = 100
	friction           = 0.7

	cellSize = 5

	gravityStrength = 1
)

type Game struct {
	cells          []*Cell
	rules          [10][10]float64
	requiresInput  bool
	partitionBoard [10][10][]*Cell
}

func (g *Game) CalculateForce(distance float64, attractionForce float64) float64 {
	beta := 0.3
	if distance < beta {
		return distance/beta - 1
	} else if beta < distance && distance < 1 {
		return attractionForce * (1 - math.Abs(2*distance-1-beta)/(1-beta))
	} else {
		return 0
	}
}

func (g *Game) PrintRules() {
	for row := 0; row < 4; row++ {
		for column := 0; column < 4; column++ {
			fmt.Print(g.rules[row][column], " ")
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func (g *Game) RandomizeRules() {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			g.rules[i][j] = rand.Float64()*2 - 1
		}
	}
	g.PrintRules()
}

func (g *Game) CalculatePartialForce(a *Cell, cells []*Cell) Vector2 {
	partialForce := Vector2{0, 0}
	for j := 0; j < len(cells); j++ {
		if a == cells[j] {
			continue
		}

		direction := Vector2{cells[j].position.x - a.position.x, cells[j].position.y - a.position.y}
		distance := direction.Magnitude()

		if distance <= 0.0001 {
			//fmt.Printf("PARTICLES TOO CLOSE: %p AND %p\n", a, cells[j])
			continue
		}

		if distance < attractionDistance {
			force := g.CalculateForce(distance/attractionDistance, g.rules[a.cellType][cells[j].cellType])
			direction.Normalize()
			partialForce.x += direction.x * force
			partialForce.y += direction.y * force
		}
	}
	return partialForce
}

func ConvertToPartitionBoardIndex(position float64) uint {
	index := uint(position) / attractionDistance

	if index == screenHeight/attractionDistance {
		index -= 1
	}

	return index
}

func (g *Game) Update() error {
	if (g.requiresInput && inpututil.IsKeyJustPressed(ebiten.KeySpace)) || !g.requiresInput {
		/*
			for i := 0; i < len(g.cells); i++ {
				totalForce := Vector2{0, 0}
				a := g.cells[i]
				for j := 0; j < len(g.cells); j++ {
					if i == j {
						continue
					}

					b := g.cells[j]

					direction := Vector2{b.position.x - a.position.x, b.position.y - a.position.y}
					distance := direction.Magnitude()

					if distance < attractionDistance {
						force := g.CalculateForce(distance/attractionDistance, g.rules[a.cellType][b.cellType])
						direction.Normalize()
						totalForce.x += direction.x * force
						totalForce.y += direction.y * force
					}
				}

				totalForce.x *= attractionDistance
				totalForce.y *= attractionDistance

				g.cells[i].velocity.x *= friction
				g.cells[i].velocity.y *= friction

				g.cells[i].velocity.x += totalForce.x * deltaT
				g.cells[i].velocity.y += totalForce.y * deltaT
			}
		*/

		//fmt.Println("NEW UPDATE")

		for i := 0; i < len(g.cells); i++ {
			totalForce := Vector2{0, 0}

			currentXIndex := ConvertToPartitionBoardIndex(g.cells[i].position.x)
			currentYIndex := ConvertToPartitionBoardIndex(g.cells[i].position.y)

			// top left
			if (currentXIndex >= 1) && (currentYIndex >= 1) {
				partialForce := g.CalculatePartialForce(g.cells[i], g.partitionBoard[currentYIndex-1][currentXIndex-1])

				totalForce.x += partialForce.x
				totalForce.y += partialForce.y
			}

			// top middle
			if currentYIndex >= 1 {
				partialForce := g.CalculatePartialForce(g.cells[i], g.partitionBoard[currentYIndex-1][currentXIndex])

				totalForce.x += partialForce.x
				totalForce.y += partialForce.y
			}

			// top right
			if currentXIndex < (screenWidth/attractionDistance-1) && currentYIndex >= 1 {
				partialForce := g.CalculatePartialForce(g.cells[i], g.partitionBoard[currentYIndex-1][currentXIndex+1])

				totalForce.x += partialForce.x
				totalForce.y += partialForce.y
			}

			// middle left
			if currentXIndex >= 1 {
				partialForce := g.CalculatePartialForce(g.cells[i], g.partitionBoard[currentYIndex][currentXIndex-1])

				totalForce.x += partialForce.x
				totalForce.y += partialForce.y
			}

			// middle middle
			if true {
				partialForce := g.CalculatePartialForce(g.cells[i], g.partitionBoard[currentYIndex][currentXIndex])

				totalForce.x += partialForce.x
				totalForce.y += partialForce.y
			}

			// middle right
			if currentXIndex < (screenWidth/attractionDistance - 1) {
				partialForce := g.CalculatePartialForce(g.cells[i], g.partitionBoard[currentYIndex][currentXIndex+1])

				totalForce.x += partialForce.x
				totalForce.y += partialForce.y
			}

			// bottom left
			if currentXIndex >= 1 && currentYIndex < (screenHeight/attractionDistance-1) {
				partialForce := g.CalculatePartialForce(g.cells[i], g.partitionBoard[currentYIndex+1][currentXIndex-1])

				totalForce.x += partialForce.x
				totalForce.y += partialForce.y
			}

			// bottom middle
			if currentYIndex < (screenHeight/attractionDistance - 1) {
				partialForce := g.CalculatePartialForce(g.cells[i], g.partitionBoard[currentYIndex+1][currentXIndex])

				totalForce.x += partialForce.x
				totalForce.y += partialForce.y
			}

			// bottom right
			if currentYIndex < (screenHeight/attractionDistance-1) && currentXIndex < (screenWidth/attractionDistance-1) {
				partialForce := g.CalculatePartialForce(g.cells[i], g.partitionBoard[currentYIndex+1][currentXIndex+1])

				totalForce.x += partialForce.x
				totalForce.y += partialForce.y
			}

			totalForce.x *= attractionDistance
			totalForce.y *= attractionDistance

			/*
				if totalForce.x >= 1000000 || math.IsNaN(totalForce.x) {
					totalForce.x = 100000
				}

				if totalForce.y >= 1000000 || math.IsNaN(totalForce.y) {
					totalForce.y = 100000
				}
			*/

			g.cells[i].velocity.x *= friction
			g.cells[i].velocity.y *= friction

			g.cells[i].velocity.x += totalForce.x * deltaT
			g.cells[i].velocity.y += totalForce.y * deltaT

			//fmt.Println()
			//fmt.Println(totalForce.x)
			//fmt.Println(g.cells[i].velocity.x)
		}

		for i := 0; i < len(g.cells); i++ {

			oldXIndex := ConvertToPartitionBoardIndex(g.cells[i].position.x)
			oldYIndex := ConvertToPartitionBoardIndex(g.cells[i].position.y)

			//println("")
			//println(c.position.x)

			g.cells[i].position.x += g.cells[i].velocity.x * deltaT
			g.cells[i].position.y += g.cells[i].velocity.y * deltaT

			//println(c.position.x)

			/*
				if g.cells[i].position.x < 0 {
					g.cells[i].position.x += float64(screenWidth)
				}

				if g.cells[i].position.y < 0 {
					g.cells[i].position.y += float64(screenHeight)
				}

				if g.cells[i].position.x > screenWidth {
					g.cells[i].position.x = float64(uint(g.cells[i].position.x) % screenWidth)
				}

				if g.cells[i].position.y > screenHeight {
					g.cells[i].position.y = float64(uint(g.cells[i].position.y) % screenHeight)
				}
			*/

			if g.cells[i].position.x < 0 {
				g.cells[i].position.x = screenWidth / 2
			}

			if g.cells[i].position.y < 0 {
				g.cells[i].position.y = screenHeight / 2
			}

			if g.cells[i].position.x > screenWidth {
				g.cells[i].position.x = screenWidth / 2
			}

			if g.cells[i].position.y > screenHeight {
				g.cells[i].position.y = screenHeight / 2
			}

			newXIndex := ConvertToPartitionBoardIndex(g.cells[i].position.x)
			newYIndex := ConvertToPartitionBoardIndex(g.cells[i].position.y)

			//fmt.Printf("%d %d %d %d\n", g.cells[i].position.x, g.cells[i].position.y, newXIndex, newYIndex)

			j := 0

			//println("NEWNEWNEWNEWNEWNEWNENWENWNEWNENWENWNEWNENWENWENWN")
			//println(g.cells[i])
			for o := 0; o < len(g.partitionBoard[oldYIndex][oldXIndex]); o++ {
				//fmt.Printf("%p\n", g.partitionBoard[oldYIndex][oldXIndex][o])
			}

			for j < len(g.partitionBoard[oldYIndex][oldXIndex]) {
				if g.partitionBoard[oldYIndex][oldXIndex][j] == g.cells[i] {
					//fmt.Printf("DELETING:\n")
					//fmt.Printf("%p\n", g.partitionBoard[oldYIndex][oldXIndex][j])
					//fmt.Printf("%p\n\n", g.cells[i])
					g.partitionBoard[oldYIndex][oldXIndex] = slices.Delete(g.partitionBoard[oldYIndex][oldXIndex], j, j+1)
					//g.partitionBoard[oldYIndex][oldXIndex][j] = g.partitionBoard[oldYIndex][oldXIndex][len(g.partitionBoard[oldYIndex][oldXIndex])-1]
					//g.partitionBoard[oldYIndex][oldXIndex] = g.partitionBoard[oldYIndex][oldXIndex][:len(g.partitionBoard[oldYIndex][oldXIndex])-1]
					//fmt.Println("DELETING")
					break
				}
				j++
			}

			//println("------------------------------------------------------------")

			for o := 0; o < len(g.partitionBoard[oldYIndex][oldXIndex]); o++ {
				//fmt.Printf("%p\n", g.partitionBoard[oldYIndex][oldXIndex][o])
			}

			//println("NEWNEWNEWNEWNEWNEWNENWENWNEWNENWENWNEWNENWENWENWN")

			for o := 0; o < len(g.partitionBoard[newYIndex][newXIndex]); o++ {
				//fmt.Printf("%p\n", g.partitionBoard[newYIndex][newXIndex][o])
			}
			//println("------------------------------------------------------------")
			g.partitionBoard[newYIndex][newXIndex] = append(g.partitionBoard[newYIndex][newXIndex], g.cells[i])

			for o := 0; o < len(g.partitionBoard[newYIndex][newXIndex]); o++ {
				//fmt.Printf("%p\n", g.partitionBoard[newYIndex][newXIndex][o])
			}

		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.requiresInput = !g.requiresInput
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.RandomizeRules()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for i := 0; i < len(g.cells); i++ {
		g.cells[i].Draw(screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) AddRule(a CellType, b CellType, force float64) {
	g.rules[a][b] = force
}

func CellConstructor(cellType CellType, position Vector2, size uint8) Cell {
	var cellColor color.Color
	switch cellType {
	case RedCell:
		cellColor = color.RGBA{255, 0, 0, 255}
	case BlueCell:
		cellColor = color.RGBA{0, 0, 255, 255}
	case GreenCell:
		cellColor = color.RGBA{0, 255, 0, 255}
	case WhiteCell:
		cellColor = color.RGBA{255, 255, 255, 255}
	}
	return Cell{1.0, cellType, cellColor, Vector2{0, 0}, position, size}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")

	g := &Game{requiresInput: true}

	i := 0
	for y := 0; y < 1000; y++ {
		var cellType CellType
		switch i % 4 {
		case 0:
			cellType = RedCell
		case 1:
			cellType = BlueCell
		case 2:
			cellType = GreenCell
		case 3:
			cellType = WhiteCell
		}

		randX := rand.Intn(screenWidth - cellSize)
		randY := rand.Intn(screenHeight - cellSize)

		c := CellConstructor(cellType, Vector2{float64(randX), float64(randY)}, 2)
		//g.cells[i] = c
		g.cells = append(g.cells, &c)

		xIndex := ConvertToPartitionBoardIndex(c.position.x)
		yIndex := ConvertToPartitionBoardIndex(c.position.y)

		g.partitionBoard[yIndex][xIndex] = append(g.partitionBoard[yIndex][xIndex], &c)
		i++
	}

	//fmt.Printf("%p\n", &g.partitionBoard[0][0][0])
	//fmt.Printf("%p\n", &g.cells[0])

	/*
		g.AddRule(RedCell, RedCell, -1)
		g.AddRule(RedCell, GreenCell, 1)
		g.AddRule(RedCell, BlueCell, 1)
		g.AddRule(RedCell, WhiteCell, 1)

		g.AddRule(GreenCell, RedCell, 1)
		g.AddRule(GreenCell, GreenCell, -1)
		g.AddRule(GreenCell, BlueCell, 1)
		g.AddRule(GreenCell, WhiteCell, -1)

		g.AddRule(BlueCell, RedCell, -1)
		g.AddRule(BlueCell, GreenCell, -1)
		g.AddRule(BlueCell, BlueCell, -1)
		g.AddRule(BlueCell, WhiteCell, -1)

		g.AddRule(WhiteCell, RedCell, -1)
		g.AddRule(WhiteCell, GreenCell, 1)
		g.AddRule(WhiteCell, BlueCell, -1)
		g.AddRule(WhiteCell, WhiteCell, 1)
	*/

	g.RandomizeRules()

	println(len(g.cells))

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
