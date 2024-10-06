package main

import (
	"log"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

type TerminalColor int

const (
	TrueColor TerminalColor = iota
	Color256
	AnsiColor
)

type Terminal struct {
	fileDescriptor int
	oldState       unix.Termios
	currentState   unix.Termios
	colorSupport   TerminalColor
}

func (t *Terminal) Init() {
	t.colorSupport = t.GetBestColorSupport()
	t.EnableRawMode()
	t.HideCursor()
	t.EnableAlternateBuffer()
	t.ClearAlternateBuffer()
}

func (t *Terminal) GetColorSupport() TerminalColor {
	return t.colorSupport
}

func (t *Terminal) GetBestColorSupport() TerminalColor {
	if t.SupportsTrueColor() {
		return TrueColor
	} else if t.Supports256Color() {
		return Color256
	} else {
		return AnsiColor
	}
}

func (t *Terminal) SupportsTrueColor() bool {
	colorTerm := os.Getenv("COLORTERM")
	return colorTerm == "truecolor" || colorTerm == "24bit"
}

func (t *Terminal) Supports256Color() bool {
	envVar := os.Getenv("TERM")
	return strings.Contains(envVar, "256color")
}

func (t *Terminal) HideCursor() {
	os.Stdout.Write([]byte(escHideCursor))
}

func (t *Terminal) ShowCursor() {
	os.Stdout.Write([]byte(escShowCursor))
}

func (t *Terminal) GetTerminalSize() (int, int) {
	ws, err := unix.IoctlGetWinsize(t.fileDescriptor, unix.TIOCGWINSZ)

	if err != nil {
		log.Fatalf("Could not access terminal size: %s", err.Error())
	}

	return int(ws.Col), int(ws.Row)
}

func (t *Terminal) ApplyState(state *unix.Termios) {
	unix.IoctlSetTermios(t.fileDescriptor, unix.TCSETS, state)
}

func (t *Terminal) Restore() {
	t.ApplyState(&t.oldState)
	t.DisableAlternateBuffer()
	t.ShowCursor()
}

func (t *Terminal) EnableRawMode() {
	// Disables echo and canonical mode
	t.currentState.Lflag &^= unix.ICANON | unix.ECHO | unix.ISIG

	// Defines the minimum number of bytes before read returns
	t.currentState.Cc[unix.VMIN] = 1

	// Defines a timeout
	t.currentState.Cc[unix.VTIME] = 0

	t.ApplyState(&t.currentState)
}

func (t *Terminal) GetFd() int {
	return t.fileDescriptor
}

func (t *Terminal) EnableAlternateBuffer() {
	os.Stdout.Write([]byte(escAlternateBuff))
}

func (t *Terminal) DisableAlternateBuffer() {
	os.Stdout.Write([]byte(escExitAlternate))
}

func (t *Terminal) ClearAlternateBuffer() {
	os.Stdout.Write([]byte(escClearScreen + escMoveCursorTop))
}

func NewTerminal() *Terminal {
	fd := int(os.Stdin.Fd())
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)

	if err != nil {
		log.Fatalf("Standard input is not a terminal: '%s'", err.Error())
	}

	return &Terminal{
		fileDescriptor: fd,
		oldState:       *termios,
		currentState:   *termios,
		colorSupport:   AnsiColor,
	}
}
