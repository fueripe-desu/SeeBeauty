package main

import (
	"strings"
)

type WordInfo struct {
	word   string
	length int
}

type Text struct {
	Text       string
	Position   *Position
	Dimensions *Dimensions
	Padding    *Padding
	Border     *Border
	Props      *TextProps
}

func (t *Text) Render() (*Matrix, int, int) {
	data := newTextData(t.Text, t.Position, t.Dimensions, t.Padding, t.Border, t.Props)
	x, y := data.getPosition()
	textW, textH := t.calculateTextbox(data)
	textMatrix := t.createTextMatrix(textW, textH, data)
	spacedMatrix := t.calculateSpacing(textMatrix, data)

	if data.hasBorder() {
		t.placeBorder(spacedMatrix, data)
	}

	return spacedMatrix, x, y
}

func (t *Text) calculateSpacing(matrix *Matrix, data *textData) *Matrix {
	// If one of the sizes has a border, then it evals to 1, otherwise 0
	borderSize := BoolToInt(data.hasBorder())
	pt, pr, pb, pl := data.getPadding()

	width := borderSize + pl + matrix.Width() + pr + borderSize
	height := borderSize + pt + matrix.Height() + pb + borderSize

	newMatrix := NewMatrix(width, height)
	newMatrix.ForEach(func(colIndex int, rowIndex int, element rune, end bool) rune {
		return rune('0')
	})
	newMatrix.PlaceMatrix(borderSize+pl+1, borderSize+pt+1, matrix)
	return newMatrix
}

func (t *Text) calculateTextbox(data *textData) (int, int) {
	width, height := data.getDimensions()
	pt, pr, pb, pl := data.getPadding()
	borderSize := BoolToInt(data.hasBorder())

	if width > 0 {
		width = width - pr - pl - borderSize - borderSize

		// If the final fixed width is less than 1, than it will default to 1
		if width < 1 {
			width = 1
		}
	}

	if height > 0 {
		height = height - pt - pb - borderSize - borderSize

		// If the final fixed height is less than 1, than it will default to 1
		if height < 1 {
			height = 1
		}
	}

	return width, height
}

func (t *Text) createTextMatrix(textW int, textH int, data *textData) *Matrix {
	maxLines, ellipsis, wordWrap := data.getProps()

	fixedH := textH > 0
	fixedW := textW > 0

	// If height is fixed but width is auto, the textbox will scale horizontally
	// making the text single-line
	if !fixedW && fixedH {
		return t.placeSingleline(textH)
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
		if textW < 4 {
			newEllipsis = false
		}
		if wordWrap {
			return t.placeMultilineWrap(textW, maxLines, newEllipsis)
		} else {
			return t.placeMultilineBW(textW, maxLines, newEllipsis)
		}
	}

	// If both width and height are fixed, the textbox will not grow and the text
	// will be cropped if it does not fit inside the box
	if fixedH && fixedW {
		newMaxLines := maxLines
		newEllipsis := ellipsis

		// If max lines is not specified, it becomes equal to the height
		if maxLines <= 0 {
			newMaxLines = textH
		}

		// Does not allow the ellipsis to show if the width is less than 4 characters
		if textW < 4 {
			newEllipsis = false
		}

		var result *Matrix

		if wordWrap {
			result = t.placeMultilineWrap(textW, newMaxLines, newEllipsis)
		} else {
			result = t.placeMultilineBW(textW, newMaxLines, newEllipsis)
		}

		// Enforce height on textbox
		result.GrowV(textH - result.Height())

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

func (t *Text) placeBorder(matrix *Matrix, data *textData) {
	if !data.hasBorder() {
		return
	}

	btl, bt, btr, br, bbl, bb, bbr, bl := data.getBorderChars()
	matrix.Border(1, bt, bl, bb, br, btl, btr, bbl, bbr)
}

type textData struct {
	text       string
	position   *Position
	dimensions *Dimensions
	padding    *Padding
	border     *Border
	props      *TextProps
}

func (d *textData) getText() string {
	return strings.Trim(d.text, " ")
}

func (d *textData) getPosition() (int, int) {
	if d.position != nil {
		return d.position.Eval()
	}
	return 1, 1
}

func (d *textData) getDimensions() (int, int) {
	if d.dimensions != nil {
		return d.dimensions.Eval()
	}
	return 1, 1
}

func (d *textData) getPadding() (int, int, int, int) {
	if d.padding != nil {
		return d.padding.Eval()
	}
	return 0, 0, 0, 0
}

func (d *textData) getBorders() (bool, bool, bool, bool) {
	if d.border != nil {
		return d.border.Eval()
	}
	return false, false, false, false
}

func (d *textData) getBorderSizes() (int, int, int, int) {
	if d.border != nil {
		return d.border.EvalSizes()
	}
	return 0, 0, 0, 0
}

func (d *textData) getBorderChars() (rune, rune, rune, rune, rune, rune, rune, rune) {
	if d.border != nil {
		return d.border.EvalBorderRunes()
	}
	return rune(' '), rune(' '), rune(' '), rune(' '), rune(' '), rune(' '), rune(' '), rune(' ')
}

func (d *textData) hasBorder() bool {
	bt, br, bb, bl := d.getBorderSizes()
	return bt+br+bb+bl > 0
}

func (d *textData) getProps() (int, bool, bool) {
	if d.props != nil {
		return d.props.Eval()
	}
	return 0, false, false
}

func newTextData(
	text string,
	position *Position,
	dimensions *Dimensions,
	padding *Padding,
	border *Border,
	props *TextProps,
) *textData {
	return &textData{
		text:       text,
		position:   position,
		dimensions: dimensions,
		padding:    padding,
		border:     border,
		props:      props,
	}
}
