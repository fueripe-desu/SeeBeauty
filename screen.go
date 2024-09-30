package main

type Screen interface {
	OnEvent(ctx *Context, event Event)
	Update(ctx *Context)
	View(ctx *Context) Component
}
