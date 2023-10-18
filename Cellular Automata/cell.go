package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CellType uint8

const (
	GreenCell CellType = iota
	RedCell
	BlueCell
	WhiteCell
)

type Cell struct {
	strength  float64
	cellType  CellType
	cellColor color.Color
	velocity  Vector2
	position  Vector2
	size      uint8
}

func (c *Cell) CalculateNewVelocity(g *Game) {
	for i := 0; i < len(g.cells); i++ {
		if *g.cells[i] == *c {
			continue
		}
		direction := Vector2{g.cells[i].position.x - c.position.x, g.cells[i].position.y - c.position.y}
		squareDistance := math.Pow(Distance(c.position, g.cells[i].position), 2)
		attractionForce := g.rules[c.cellType][g.cells[i].cellType]
		c.velocity.x += direction.x / squareDistance * attractionForce
		c.velocity.y += direction.y / squareDistance * attractionForce
	}
}

func (c *Cell) Update(g *Game) {
	c.position.x += c.velocity.x * deltaT
	c.position.y += c.velocity.y * deltaT

	if c.position.x < 0 {
		c.position.x = 0
		c.velocity.x *= -1
	}
	if c.position.y < 0 {
		c.position.y = 0
		c.velocity.y *= -1
	}
	if c.position.x+float64(c.size) > screenWidth {
		c.position.x = float64(screenWidth - uint(c.size))
		c.velocity.x *= -1
	}
	if c.position.y+float64(c.size) > screenHeight {
		c.position.y = float64(screenHeight - uint(c.size))
		c.velocity.y *= -1
	}

	//c.velocity.x *= 0.9
	//c.velocity.y *= 0.9
}

func (c *Cell) Draw(screen *ebiten.Image) {
	//vector.DrawFilledRect(screen, float32(c.position.x), float32(c.position.y), float32(c.size), float32(c.size), c.cellColor, false)
	vector.DrawFilledCircle(screen, float32(c.position.x), float32(c.position.y), float32(c.size), c.cellColor, true)
}
