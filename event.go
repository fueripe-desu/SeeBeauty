package main

type Event interface {
	Payload() map[string]any
}

type OnWindowCreate struct {
}

func (e *OnWindowCreate) Payload() map[string]any {
	return nil
}

type OnCreate struct {
}

func (e *OnCreate) Payload() map[string]any {
	return nil
}
