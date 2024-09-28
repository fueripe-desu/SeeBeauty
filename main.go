package main

import (
	"fmt"
	"os"
)

type MainScreen struct {
}

func (s *MainScreen) OnEvent(ctx *Context, event Event) {
	switch event.(type) {
	case *OnWindowCreate:
		fmt.Println("Window created!")
	case *OnCreate:
		fmt.Println("Screen created!")
	}
}

func (s *MainScreen) Update(ctx *Context) {
	var b = make([]byte, 1)
	os.Stdin.Read(b)
	fmt.Println("Something")
	ctx.SendSignal(SigExit)
}

func main() {
	term := NewTerminal()
	term.Init()

	renderer := NewRenderer(term)
	mainScreen := &MainScreen{}

	renderer.OpenScreen(mainScreen)
}
