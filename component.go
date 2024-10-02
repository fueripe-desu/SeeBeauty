package main

type Component interface {
	Render() *Matrix
}

type Text struct {
	PosX  int
	PosY  int
	Value string
}

func (t *Text) Render() *Matrix {
	w := len(t.Value)
	matrix := NewMatrix(w, 1)

	for i, r := range t.Value {
		matrix.Place(i+1, 1, r)
	}

	return matrix
}

func NewText(x int, y int, value string) *Text {
	return &Text{
		PosX:  x,
		PosY:  y,
		Value: value,
	}
}
