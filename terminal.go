package main

import (
	"log"
	"os"

	"golang.org/x/sys/unix"
)

type Terminal struct {
	fileDescriptor int
	oldState       unix.Termios
	currentState   unix.Termios
}

func (t *Terminal) Init() {
	t.EnableRawMode()
	t.HideCursor()
	t.EnableAlternateBuffer()
	t.ClearAlternateBuffer()
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
	}
}
