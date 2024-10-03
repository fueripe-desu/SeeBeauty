package main

type Component interface {
	Render() (*Matrix, int, int)
}

type Text struct {
	Text string

	PosX   *Integer
	PosY   *Integer
	Border *Bool
}

func (t *Text) Render() (*Matrix, int, int) {
	// Params Evaluation
	x, y := t.evalPosition()
	border := t.evalBorder()

	w := len(t.Text)
	matrix := t.generateMatrix(w, border)
	matrix.ForEach(func(colIndex int, rowIndex int, element rune, end bool) rune {
		return rune('0')
	})
	t.placeBorder(matrix, border)

	// for i, r := range t.Text {
	// matrix.Place(i+1, 1, r)
	// }

	return matrix, x, y
}

func (t *Text) evalBorder() bool {
	if t.Border != nil {
		return t.Border.val
	}

	return false
}

func (t *Text) placeBorder(matrix *Matrix, border bool) {
	if !border {
		return
	}

	matrix.Border(
		1,
		rune('─'),
		rune('│'),
		rune('─'),
		rune('│'),
		rune('┌'),
		rune('┐'),
		rune('└'),
		rune('┘'),
	)
}

func (t *Text) generateMatrix(width int, border bool) *Matrix {
	if border {
		return NewMatrix(width+4, 1+4)
	} else {
		return NewMatrix(width, 1)
	}
}

func (t *Text) evalPosition() (int, int) {
	x := 1
	y := 1

	if t.PosX != nil {
		x = t.PosX.val
	}

	if t.PosY != nil {
		y = t.PosY.val
	}

	return x, y
}
