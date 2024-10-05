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
		Text:          "Hello world!",
		PosX:          NewInt(40),
		PosY:          NewInt(10),
		Border:        NewBool(true),
		PaddingTop:    NewInt(2),
		PaddingBottom: NewInt(2),
		PaddingLeft:   NewInt(2),
		PaddingRight:  NewInt(2),
		Width:         NewInt(18),
		Height:        NewInt(15),
		EllipsisWrap:  NewBool(false),
	}
}

func main() {
	term := NewTerminal()
	term.Init()

	renderer := NewRenderer(term)
	mainScreen := &MainScreen{}

	renderer.OpenScreen(mainScreen)
}
