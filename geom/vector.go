package geom

import "fmt"

// not actually used in the project, just used to help create test templating

type Vector2d struct {
	X float64
	Y float64
}

func (a Vector2d) Add(b Vector2d) Vector2d {
	return Vector2d{a.X + b.X, a.Y + b.Y}
}

func AddVector2ds(vecA Vector2d, vecB Vector2d) Vector2d {
	return Vector2d{vecA.X + vecB.X, vecA.Y + vecB.Y}
}

func TestAddVectors() {
	vecA := Vector2d{1.0, 2.0}
	vecB := Vector2d{2.0, 1.0}
	vecC := vecA.Add(vecB)
	fmt.Println(vecC)
}
