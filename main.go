package main

import (
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
	return NewText(0, 0, "Hello world!")
}

func main() {
	term := NewTerminal()
	term.Init()

	renderer := NewRenderer(term)
	mainScreen := &MainScreen{}

	renderer.OpenScreen(mainScreen)
}
