package main

type Component interface {
	Render() (*Matrix, int, int)
}
