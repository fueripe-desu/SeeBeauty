package main

type Context struct {
	signals Queue[Signal]
	refresh bool
	window  *WindowParams
}

func (c *Context) SendSignal(signal Signal) {
	c.signals.Enqueue(signal)
}

func NewContext(width int, height int) *Context {
	return &Context{
		signals: Queue[Signal]{},
		refresh: true,
		window: &WindowParams{
			Width:  width,
			Height: height,
		},
	}
}

type WindowParams struct {
	Width  int
	Height int
}
