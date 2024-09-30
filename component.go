package main

type Component interface {
	Render() string
}

type Text struct {
	PosX  int
	PosY  int
	Value string
}

func (t *Text) Render() string {
	return t.Value
}

func NewText(x int, y int, value string) *Text {
	return &Text{
		PosX:  x,
		PosY:  y,
		Value: value,
	}
}
