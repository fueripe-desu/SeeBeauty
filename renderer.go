package main

import (
	"fmt"
	"log"
	"os"
)

type Renderer struct {
	terminal *Terminal
	context  *Context
}

func (r *Renderer) OpenScreen(screen Screen) {
	r.checkContext()

	if screen == nil {
		log.Fatal("Cannot open screen: screen is nil.")
	}

	screen.OnEvent(r.context, &OnWindowCreate{})
	screen.OnEvent(r.context, &OnCreate{})

	// Main loop
	for {
		if r.context.refresh {
			fmt.Print(screen.View(r.context).Render())
			fmt.Print(r.terminal.GetTerminalSize())
			r.context.refresh = false
		}

		for !(r.context.signals.IsEmpty()) {
			r.handleSignal(r.context.signals.Dequeue())
		}
		screen.Update(r.context)
	}
}

func (r *Renderer) handleSignal(signal Signal) {
	switch signal {
	case SigExit:
		r.terminal.Restore()
		os.Exit(0)
	}
}

func (r *Renderer) checkContext() {
	if r.context == nil {
		log.Fatal("Cannot initialize renderer if context is nil.")
	}
}

func NewRenderer(term *Terminal) *Renderer {
	ctx := NewContext()

	return &Renderer{
		terminal: term,
		context:  ctx,
	}
}
