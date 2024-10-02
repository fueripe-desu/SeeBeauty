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
	return &Text{
		Text: "Hello world!",
		PosX: NewOption(40),
		PosY: NewOption(10),
	}
}

func main() {
	term := NewTerminal()
	term.Init()

	renderer := NewRenderer(term)
	mainScreen := &MainScreen{}

	renderer.OpenScreen(mainScreen)
}
