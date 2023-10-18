package main

import "math"

type Vector2 struct {
	x float64
	y float64
}

func Distance(a Vector2, b Vector2) float64 {
	diffVector := Vector2{b.x - a.x, b.y - a.y}
	return diffVector.Magnitude()
}

func (v *Vector2) Magnitude() float64 {
	return math.Sqrt(math.Pow(v.x, 2) + math.Pow(v.y, 2))
}

func (v *Vector2) Normalize() {
	magnitude := v.Magnitude()
	v.x /= magnitude
	v.y /= magnitude
}
