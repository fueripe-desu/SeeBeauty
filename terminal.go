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

func (t *Terminal) ApplyState(state *unix.Termios) {
	unix.IoctlSetTermios(t.fileDescriptor, unix.TCSETS, state)
}

func (t *Terminal) Restore() {
	t.ApplyState(&t.oldState)
}

func (t *Terminal) EnableRawMode() {
	// Disables echo and canonical mode
	t.currentState.Lflag &^= unix.ECHO | unix.ICANON

	// Defines the minimum number of bytes before read returns
	t.currentState.Cc[unix.VMIN] = 1

	// Defines a timeout
	t.currentState.Cc[unix.VTIME] = 0

	t.ApplyState(&t.currentState)
}

func (t *Terminal) GetFd() int {
	return t.fileDescriptor
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
