package main

const (
	escClearScreen   = "\033[2J"
	escMoveCursorTop = "\033[H"
	escAlternateBuff = "\033[?1049h"
	escExitAlternate = "\033[?1049l"
	escHideCursor    = "\033[?25l"
	escShowCursor    = "\033[?25h"
)
