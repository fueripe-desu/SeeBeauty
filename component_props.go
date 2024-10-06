package main

type BorderType int
type BorderStyle int

const (
	BorderTop BorderType = iota
	BorderRight
	BorderBottom
	BorderLeft
)

const (
	BorderSolid BorderStyle = iota
	BorderThick
	BorderDashed
	BorderDouble
	BorderRounded
)

type Position struct {
	x int
	y int
}

func (p *Position) Eval() (int, int) {
	return defaultToOne(p.x), defaultToOne(p.y)
}

func NewXY(x int, y int) *Position {
	return &Position{x: x, y: y}
}

type Dimensions struct {
	width  int
	height int
}

func (d *Dimensions) Eval() (int, int) {
	return defaultToOne(d.width), defaultToOne(d.height)
}

func NewWH(width int, height int) *Dimensions {
	return &Dimensions{width: width, height: height}
}

type Padding struct {
	t int
	r int
	b int
	l int
}

func (p *Padding) Eval() (int, int, int, int) {
	return defaultToZero(p.t), defaultToZero(p.r), defaultToZero(p.b), defaultToZero(p.l)
}

func NewPTRBL(t int, r int, b int, l int) *Padding {
	return &Padding{t: t, r: r, b: b, l: l}
}

func NewPHV(h int, v int) *Padding {
	return &Padding{t: v, r: h, b: v, l: h}
}

func NewPadding(p int) *Padding {
	return &Padding{t: p, r: p, b: p, l: p}
}

type BorderSide struct {
	borderType  BorderType
	borderStyle BorderStyle
}

func (s *BorderSide) Eval() (rune, rune, rune) {
	switch s.borderType {
	case BorderTop:
		switch s.borderStyle {
		case BorderSolid:
			return rune('┌'), rune('─'), rune('┐')
		case BorderThick:
			return rune('┏'), rune('━'), rune('┓')
		case BorderDashed:
			return rune('+'), rune('-'), rune('+')
		case BorderDouble:
			return rune('╔'), rune('═'), rune('╗')
		case BorderRounded:
			return rune('╭'), rune('─'), rune('╮')
		}
	case BorderRight:
		switch s.borderStyle {
		case BorderSolid:
			return rune('│'), rune(0), rune(0)
		case BorderThick:
			return rune('┃'), rune(0), rune(0)
		case BorderDashed:
			return rune('|'), rune(0), rune(0)
		case BorderDouble:
			return rune('║'), rune(0), rune(0)
		case BorderRounded:
			return rune('│'), rune(0), rune(0)
		}
	case BorderBottom:
		switch s.borderStyle {
		case BorderSolid:
			return rune('└'), rune('─'), rune('┘')
		case BorderThick:
			return rune('┗'), rune('━'), rune('┛')
		case BorderDashed:
			return rune('+'), rune('-'), rune('+')
		case BorderDouble:
			return rune('╚'), rune('═'), rune('╝')
		case BorderRounded:
			return rune('╰'), rune('─'), rune('╯')
		}
	case BorderLeft:
		switch s.borderStyle {
		case BorderSolid:
			return rune('│'), rune(0), rune(0)
		case BorderThick:
			return rune('┃'), rune(0), rune(0)
		case BorderDashed:
			return rune('|'), rune(0), rune(0)
		case BorderDouble:
			return rune('║'), rune(0), rune(0)
		case BorderRounded:
			return rune('│'), rune(0), rune(0)
		}
	}

	panic("Invalid border side")
}

func NewBorderSide(borderStyle BorderStyle) *BorderSide {
	return &BorderSide{borderStyle: borderStyle}
}

type Border struct {
	t *BorderSide
	r *BorderSide
	b *BorderSide
	l *BorderSide
}

func (b *Border) Eval() (bool, bool, bool, bool) {
	return b.t != nil, b.r != nil, b.b != nil, b.l != nil
}

func (b *Border) EvalSizes() (int, int, int, int) {
	bt, br, bb, bl := b.Eval()
	return BoolToInt(bt), BoolToInt(br), BoolToInt(bb), BoolToInt(bl)
}

func (b *Border) EvalBorderRunes() (rune, rune, rune, rune, rune, rune, rune, rune) {
	var btl, bt, btr, br, bbl, bb, bbr, bl rune

	if b.t != nil {
		btl, bt, btr = b.t.Eval()
	} else {
		btl, bt, btr = rune(' '), rune(' '), rune(' ')
	}

	if b.r != nil {
		br, _, _ = b.r.Eval()
	} else {
		br = rune(' ')
	}

	if b.b != nil {
		bbl, bb, bbr = b.b.Eval()
	} else {
		bbl, bb, bbr = rune(' '), rune(' '), rune(' ')
	}

	if b.l != nil {
		bl, _, _ = b.l.Eval()
	} else {
		bl = rune(' ')
	}

	return btl, bt, btr, br, bbl, bb, bbr, bl
}

func NewBTRBL(t *BorderSide, r *BorderSide, b *BorderSide, l *BorderSide) *Border {
	if t != nil {
		t.borderType = BorderTop
	}

	if r != nil {
		r.borderType = BorderRight
	}

	if b != nil {
		b.borderType = BorderBottom
	}

	if l != nil {
		l.borderType = BorderLeft
	}

	return &Border{t: t, r: r, b: b, l: l}
}

func NewBHV(h *BorderSide, v *BorderSide) *Border {
	var t, r, b, l *BorderSide

	if h != nil {
		right := *h
		left := *h

		right.borderType = BorderRight
		left.borderType = BorderLeft

		r = &right
		l = &left
	} else {
		r = nil
		l = nil
	}

	if v != nil {
		top := *v
		bottom := *v

		top.borderType = BorderTop
		bottom.borderType = BorderBottom

		t = &top
		b = &bottom
	} else {
		t = nil
		b = nil
	}

	return &Border{t: t, r: r, b: b, l: l}
}

func NewBorder(s *BorderSide) *Border {
	var t, r, b, l *BorderSide

	if s != nil {
		top := *s
		right := *s
		bottom := *s
		left := *s

		top.borderType = BorderTop
		right.borderType = BorderRight
		bottom.borderType = BorderBottom
		left.borderType = BorderLeft

		t = &top
		r = &right
		b = &bottom
		l = &left
	} else {
		t = nil
		r = nil
		b = nil
		l = nil
	}

	return &Border{t: t, r: r, b: b, l: l}
}

type TextProps struct {
	maxLines int
	ellipsis bool
	wordWrap bool
}

func (p *TextProps) Eval() (int, bool, bool) {
	return defaultToZero(p.maxLines), p.ellipsis, p.wordWrap
}

func NewTextProps(maxLines int, ellipsis bool, wordWrap bool) *TextProps {
	return &TextProps{maxLines: maxLines, ellipsis: ellipsis, wordWrap: wordWrap}
}

func defaultToOne(a int) int {
	if a <= 0 {
		return 1
	}
	return a
}

func defaultToZero(a int) int {
	if a < 0 {
		return 0
	}
	return a
}
