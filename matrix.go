package main

import (
	"log"
	"strings"
)

type ElementCallback = func(element rune) rune
type RowCallback = func(row []rune)
type RevolutionCallback = func(index int, element rune) rune

type Matrix struct {
	data   [][]rune
	width  int
	height int
}

func (m *Matrix) Disnulify() {
	m.ForEach(func(element rune) rune {
		if element == rune(0) {
			return rune(' ')
		}

		return element
	}, nil)
}

func (m *Matrix) Clear() {
	m.ForEach(func(element rune) rune {
		return rune(' ')
	}, nil)
}

func (m *Matrix) GrowV(n int) {
	if n <= 0 {
		return
	}

	m.height += n
	for i := 0; i < n; i++ {
		row := make([]rune, m.width)
		for i := range row {
			row[i] = rune(' ')
		}

		m.data = append(m.data, row)
	}
}

func (m *Matrix) GrowH(n int) {
	if n <= 0 {
		return
	}

	m.width += n
	for outer := range m.data {
		row := m.data[outer]
		size := len(row)

		for size < m.width {
			row = append(row, rune(' '))
			size = len(row)
		}

		m.data[outer] = row
	}
}

func (m *Matrix) Border(
	depth int,
	t rune,
	l rune,
	b rune,
	r rune,
	tl rune,
	tr rune,
	bl rune,
	br rune,
) {
	if depth < 1 {
		log.Fatal("Depth of a revolution cannot be less than 1.")
	}

	top := depth
	bottom := m.height - (depth - 1)
	left := depth
	right := m.width - (depth - 1)

	for y := 1; y <= m.height; y++ {
		for x := 1; x <= m.width; x++ {
			if y == top {
				if x == left {
					m.Place(x, y, tl)
				} else if x == right {
					m.Place(x, y, tr)
				} else {
					m.Place(x, y, t)
				}
			} else if y == bottom {
				if x == left {
					m.Place(x, y, bl)
				} else if x == right {
					m.Place(x, y, br)
				} else {
					m.Place(x, y, b)
				}
			} else {
				if x == left {
					m.Place(x, y, l)
				} else if x == right {
					m.Place(x, y, r)
				}
			}
		}
	}
}

func (m *Matrix) Slice(x int, y int, width int, height int) *Matrix {
	if x+width > m.width || y+height > m.height || x < 0 || y < 0 {
		log.Fatal("Slice out of bounds.")
	}

	if width <= 0 || height <= 0 {
		log.Fatal("Both width and height must be greater than 0.")
	}

	matrix := NewMatrix(width, height)

	for matrixY := 1; matrixY <= height; matrixY++ {
		for matrixX := 1; matrixX <= width; matrixX++ {
			matrix.Place(matrixX, matrixY, m.data[y+matrixY][x+matrixX])
		}
	}

	return matrix
}

func (m *Matrix) PlaceMatrix(x int, y int, matrix *Matrix) {
	elementX := -1
	elementY := 0

	matrix.ForEach(
		func(element rune) rune {
			elementX++
			m.Place(x+elementX, y+elementY, element)
			return element
		},

		func(row []rune) {
			elementX = -1
			elementY++
		},
	)
}

func (m *Matrix) Place(x int, y int, element rune) {
	if x < 0 {
		log.Fatal("Cannot place an element at a negative X axis.")
	}

	if y < 0 {
		log.Fatal("Cannot place an element at a negative Y axis.")
	}

	if x > m.width {
		diff := x - m.width
		m.GrowH(diff)
	}

	if y > m.height {
		diff := y - m.height
		m.GrowV(diff)
	}

	newX := x - 1
	newY := y - 1

	m.data[newY][newX] = element
}

func (m *Matrix) ForEach(callback ElementCallback, rowCallback RowCallback) {
	if callback == nil {
		return
	}

	for outer := range m.data {
		inner := m.data[outer]
		for i := range inner {
			inner[i] = callback(inner[i])
		}

		if rowCallback != nil {
			rowCallback(inner)
		}
	}
}

func (m *Matrix) ToBuffer() string {
	var builder strings.Builder

	elementCallback := func(element rune) rune {
		builder.WriteRune(element)
		return element
	}

	rowCallback := func(row []rune) {
		builder.WriteString("\n")
	}

	m.ForEach(elementCallback, rowCallback)

	return builder.String()[:builder.Len()-1]
}

func NewMatrix(width int, height int) *Matrix {
	if width < 1 {
		log.Fatal("Matrix width should be at least 1.")
	}

	if height < 1 {
		log.Fatal("Matrix height should be at least 1.")
	}

	matrix := make([][]rune, height)

	for i := range matrix {
		matrix[i] = make([]rune, width)

		for e := range matrix[i] {
			matrix[i][e] = rune(' ')
		}
	}

	return &Matrix{
		data:   matrix,
		width:  width,
		height: height,
	}
}
