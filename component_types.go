package main

type Integer struct {
	val int
}

func NewInt(value int) *Integer {
	return &Integer{val: value}
}

type Bool struct {
	val bool
}

func NewBool(value bool) *Bool {
	return &Bool{val: value}
}
