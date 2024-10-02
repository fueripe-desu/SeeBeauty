package main

import (
	"math/rand"
	"os"
)

type MainScreen struct {
}

func (s *MainScreen) OnEvent(ctx *Context, event Event) {
	switch event.(type) {
	case *OnWindowCreate:
	case *OnCreate:
	}
}

func (s *MainScreen) Update(ctx *Context) {
	var b = make([]byte, 1)
	os.Stdin.Read(b)
	ctx.SendSignal(SigExit)
}

func (s *MainScreen) View(ctx *Context) Component {
	matrix := NewMatrix(40, 20)

	matrix.ForEach(func(element rune) rune {
		return rune('0' + rand.Intn(10))
	}, nil)

	slice := matrix.Slice(0, 0, 22, 6)

	return NewText(0, 0, slice.ToBuffer())
}

func main() {
	term := NewTerminal()
	term.Init()

	renderer := NewRenderer(term)
	mainScreen := &MainScreen{}

	renderer.OpenScreen(mainScreen)
}
