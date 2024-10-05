package main

import (
	"strings"
)

type Component interface {
	Render() (*Matrix, int, int)
}

type WordInfo struct {
	word   string
	length int
}

type Text struct {
	Text string

	// Positioning
	PosX *Integer
	PosY *Integer

	// Dimensions
	Width  *Integer
	Height *Integer

	// Padding
	PaddingTop    *Integer
	PaddingRight  *Integer
	PaddingBottom *Integer
	PaddingLeft   *Integer

	// Border
	Border *Bool

	// Text Properties
	MaxLines     *Integer
	EllipsisWrap *Bool
	WordWrap     *Bool
}

func (t *Text) Render() (*Matrix, int, int) {
	// Params Evaluation
	x, y := t.evalPosition()
	pt, pr, pb, pl := t.evalPadding()
	border := t.evalBorder()
	w, h := t.evalDimensions()
	textW, textH := t.calculateTextbox(w, h, pt, pr, pb, pl, border)
	maxLines := t.evalMaxLines()
	ellipsis := t.evalEllipsis()
	wordWrap := t.evalWordWrap()

	textMatrix := t.createTextMatrix(textW, textH, maxLines, ellipsis, wordWrap)
	spacedMatrix := t.calculateSpacing(textMatrix, pt, pr, pb, pl, border)
	t.placeBorder(spacedMatrix, border)

	return spacedMatrix, x, y
}

func (t *Text) calculateSpacing(matrix *Matrix, pt int, pr int, pb int, pl int, border bool) *Matrix {
	borderSize := 0

	if border {
		borderSize = 1
	}

	width := borderSize + pl + matrix.Width() + pr + borderSize
	height := borderSize + pt + matrix.Height() + pb + borderSize

	newMatrix := NewMatrix(width, height)
	newMatrix.PlaceMatrix(borderSize+pl+1, borderSize+pt+1, matrix)
	return newMatrix
}

func (t *Text) calculateTextbox(w int, h int, pt int, pr int, pb int, pl int, border bool) (int, int) {
	width := w
	height := h

	borderSize := 0

	if border {
		borderSize = 1
	}

	if width > 0 {
		width = w - pr - pl - borderSize - borderSize

		// If the final fixed width is less than 1, than it will default to 1
		if width < 1 {
			width = 1
		}
	}

	if height > 0 {
		height = h - pt - pb - borderSize - borderSize

		// If the final fixed height is less than 1, than it will default to 1
		if height < 1 {
			height = 1
		}
	}

	return width, height
}

func (t *Text) createTextMatrix(w int, h int, maxLines int, ellipsis bool, wordWrap bool) *Matrix {
	fixedH := h > 0
	fixedW := w > 0

	// If height is fixed but width is auto, the textbox will scale horizontally
	// making the text single-line
	if !fixedW && fixedH {
		return t.placeSingleline(h)
	}

	// If both height and width are auto, the textbox will also scale horizontally
	// making the text single-line
	if !fixedW && !fixedH {
		return t.placeSingleline(1)
	}

	// If width is fixed and height is auto, the textbox will scale vertically
	// allowing the text to have multiple lines
	if fixedW && !fixedH {
		newEllipsis := ellipsis

		// Does not allow the ellipsis to show if the width is less than 4 characters
		if w < 4 {
			newEllipsis = false
		}
		if wordWrap {
			return t.placeMultilineWrap(w, maxLines, newEllipsis)
		} else {
			return t.placeMultilineBW(w, maxLines, newEllipsis)
		}
	}

	// If both width and height are fixed, the textbox will not grow and the text
	// will be cropped if it does not fit inside the box
	if fixedH && fixedW {
		newMaxLines := maxLines
		newEllipsis := ellipsis

		// If max lines is not specified, it becomes equal to the height
		if maxLines <= 0 {
			newMaxLines = h
		}

		// Does not allow the ellipsis to show if the width is less than 4 characters
		if w < 4 {
			newEllipsis = false
		}

		var result *Matrix

		if wordWrap {
			result = t.placeMultilineWrap(w, newMaxLines, newEllipsis)
		} else {
			result = t.placeMultilineBW(w, newMaxLines, newEllipsis)
		}

		// Enforce height on textbox
		result.GrowV(h - result.Height())

		return result
	}

	panic("Unkown error")
}

func (t *Text) placeMultilineWrap(width int, maxLines int, ellipsis bool) *Matrix {
	rawWords := strings.Split(strings.Trim(t.Text, " "), " ")

	// List to store words and their lengths
	wordList := []WordInfo{}

	for _, w := range rawWords {
		if len(w) > width {
			return t.placeMultilineBW(width, maxLines, ellipsis)
		}
		wordList = append(wordList, WordInfo{word: w, length: len(w)})
	}

	matrix := NewMatrix(width, 1)

	row := ""
	line := 1

	wordCount := len(wordList)
	wordsRendered := 0

	// Iterate over the word list instead of a map
	for _, wordInfo := range wordList {
		rowSize := len(row)

		if rowSize+wordInfo.length <= width {
			row += wordInfo.word
			wordsRendered++

			// Recalculate row size
			rowSize := len(row)

			if rowSize+1 <= width {
				row += " "
			}
		} else {
			diff := width - rowSize

			// Add spaces to fill the line
			for i := 0; i < diff; i++ {
				row += " "
			}

			// Place the row into the matrix
			for i, r := range row {
				matrix.Place(i+1, line, r)
			}

			// Reset row for the next line
			row = ""
			line++

			if maxLines > 0 && line > maxLines {
				break
			}

			matrix.GrowV(1)

			// Add the current word in the row
			row += wordInfo.word
			wordsRendered++

			rowSize = len(row)

			if rowSize+1 <= width {
				row += " "
			}
		}
	}

	// If there's any leftover row, add it to the matrix
	if row != "" {
		for i, r := range row {
			matrix.Place(i+1, line, r)
		}
	}

	// Add ellipis if specified
	if ellipsis && (wordCount > wordsRendered) {
		raw := matrix.GetRow(maxLines)
		lastLine := ""

		for _, r := range raw {
			lastLine += string(r)
		}

		lastLine = strings.TrimRight(lastLine, " ")
		ellipsis := "..."

		for len(lastLine)+3 > width {
			lastLine = lastLine[:len(lastLine)-1]
		}

		lastLine += ellipsis
		finalRunes := []rune{}

		for _, r := range lastLine {
			finalRunes = append(finalRunes, r)
		}

		matrix.PlaceRow(1, maxLines, finalRunes)
	}

	return matrix
}

// Used when break word is active, this algorithm breaks the words during
// line wrapping, as opposed to just wrapping the entire word to the next
// line
func (t *Text) placeMultilineBW(width int, maxLines int, ellipsis bool) *Matrix {
	size := len(t.Text)
	matrix := NewMatrix(width, 1)

	line := 1
	index := 0
	offset := 0

	for i := 1; i <= size; i++ {
		r := rune(t.Text[index])
		x := (i - offset) % width

		if x == 1 && r == rune(' ') {
			offset++
			index++
			continue
		}

		if x == 0 {
			x = width
		}

		matrix.Place(x, line, r)

		if (i-offset)%width == 0 {
			line++

			if maxLines > 0 && line > maxLines {
				// Add ellipsis if specified
				if ellipsis && i < size {
					raw := matrix.GetRow(maxLines)
					lastLine := ""

					for _, r := range raw {
						lastLine += string(r)
					}

					lastLine = strings.TrimRight(lastLine, " ")
					ellipsis := "..."

					for len(lastLine)+3 > width {
						lastLine = lastLine[:len(lastLine)-1]
					}

					lastLine += ellipsis
					finalRunes := []rune{}

					for _, r := range lastLine {
						finalRunes = append(finalRunes, r)
					}

					matrix.PlaceRow(1, maxLines, finalRunes)
				}
				break
			}

			if i < size {
				matrix.GrowV(1)
			}
		}

		index++
	}

	return matrix
}

func (t *Text) placeSingleline(height int) *Matrix {
	size := len(t.Text)
	matrix := NewMatrix(size, height)

	for i, r := range t.Text {
		matrix.Place(i+1, 1, r)
	}

	return matrix
}

func (t *Text) evalMaxLines() int {
	maxLines := 0

	if t.MaxLines != nil {
		maxLines = t.MaxLines.val
	}

	return maxLines
}

func (t *Text) evalEllipsis() bool {
	if t.EllipsisWrap != nil {
		return t.EllipsisWrap.val
	}

	return false
}

func (t *Text) evalWordWrap() bool {
	if t.WordWrap != nil {
		return t.WordWrap.val
	}

	return false
}

func (t *Text) evalBorder() bool {
	if t.Border != nil {
		return t.Border.val
	}

	return false
}

func (t *Text) evalPadding() (int, int, int, int) {
	pt := 0
	pr := 0
	pb := 0
	pl := 0

	if t.PaddingTop != nil {
		pt = t.PaddingTop.val
	}

	if t.PaddingRight != nil {
		pr = t.PaddingRight.val
	}

	if t.PaddingBottom != nil {
		pb = t.PaddingBottom.val
	}

	if t.PaddingLeft != nil {
		pl = t.PaddingLeft.val
	}

	return pt, pr, pb, pl
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

func (t *Text) evalDimensions() (int, int) {
	w := 0
	h := 0

	if t.Width != nil {
		w = t.Width.val
	}

	if t.Height != nil {
		h = t.Height.val
	}

	return w, h
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
