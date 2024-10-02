package main

type Component interface {
	Render() (*Matrix, int, int)
}

type Option[T int] struct {
	Value    *T
	HasValue bool
}

func NewOption[T int](value T) Option[T] {
	return Option[T]{Value: &value, HasValue: true}
}

func NewNone[T int]() Option[T] {
	return Option[T]{Value: nil, HasValue: false}
}

type Text struct {
	Text string

	PosX Option[int]
	PosY Option[int]
}

func (t *Text) Render() (*Matrix, int, int) {
	w := len(t.Text)
	matrix := NewMatrix(w, 1)

	for i, r := range t.Text {
		matrix.Place(i+1, 1, r)
	}

	x, y := t.evalPosition()
	return matrix, x, y
}

func (t *Text) evalPosition() (int, int) {
	var x int
	var y int

	if t.PosX.HasValue {
		x = *(t.PosX.Value)
	} else {
		x = 1
	}

	if t.PosY.HasValue {
		y = *(t.PosY.Value)
	} else {
		y = 1
	}

	return x, y
}
