package main

import (
	"log"
	"strings"
)

type ElementCallback = func(colIndex int, rowIndex int, element rune, end bool) rune

type Matrix struct {
	data   [][]rune
	width  int
	height int
}

func (m *Matrix) Disnulify() {
	m.ForEach(func(colIndex int, rowIndex int, element rune, end bool) rune {
		if element == rune(0) {
			return rune(' ')
		}

		return element
	})
}

func (m *Matrix) Clear() {
	m.ForEach(func(colIndex int, rowIndex int, element rune, end bool) rune {
		return rune(' ')
	})
}

func (m *Matrix) Height() int {
	return m.height
}

func (m *Matrix) Width() int {
	return m.width
}

func (m *Matrix) Get(col int, row int) rune {
	if col > m.height || col < 1 {
		log.Fatal("Column out of bounds.")
	}

	if row > m.width || row < 1 {
		log.Fatal("Row out of bounds.")
	}
	return m.data[col-1][row-1]
}

func (m *Matrix) GetRow(row int) []rune {
	if row < 1 || row > m.height {
		log.Fatal("Row out of bounds.")
	}
	return m.data[row-1]

}

func (m *Matrix) GetCol(col int) []rune {
	if col < 1 || col > m.width {
		log.Fatal("Column out of bounds.")
	}

	result := []rune{}

	for i := range m.data {
		row := m.data[i]
		result = append(result, row[col-1])
	}

	return result
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
		func(colIndex int, rowIndex int, element rune, end bool) rune {
			elementX++
			m.Place(x+elementX, y+elementY, element)

			if end {
				elementX = -1
				elementY++
			}
			return element
		},
	)
}

func (m *Matrix) PlaceRow(x int, y int, row []rune) {
	if row == nil {
		log.Fatal("Cannot place a row if nil.")
	}

	if len(row) != m.width {
		log.Fatal("Row have a different width than the matrix.")
	}

	for i, r := range row {
		m.Place(x+i, y, r)
	}
}

func (m *Matrix) PlaceCol(x int, y int, col []rune) {
	if col == nil {
		log.Fatal("Cannot place a column if nil.")
	}

	if len(col) != m.height {
		log.Fatal("Column have a different height than the matrix.")
	}

	for i, r := range col {
		m.Place(x, y+i, r)
	}
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

func (m *Matrix) ForEach(callback ElementCallback) {
	if callback == nil {
		return
	}

	for outer := range m.data {
		inner := m.data[outer]
		for i := range inner {
			inner[i] = callback(
				outer,
				i,
				inner[i],
				i == (len(inner)-1),
			)
		}
	}
}

func (m *Matrix) ToBuffer() string {
	var builder strings.Builder

	m.ForEach(func(colIndex int, rowIndex int, element rune, end bool) rune {
		builder.WriteRune(element)

		if end {
			builder.WriteString("\n")
		}
		return element
	})

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
