package main

type Context struct {
	signals Queue[Signal]
}

func (c *Context) SendSignal(signal Signal) {
	c.signals.Enqueue(signal)
}

func NewContext() *Context {
	return &Context{
		signals: Queue[Signal]{},
	}
}
