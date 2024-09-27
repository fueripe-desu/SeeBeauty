package main

import (
	"fmt"
	"os"
)

func main() {
	term := NewTerminal()
	term.EnableRawMode()
	defer term.Restore()

	fmt.Println("Raw mode enabled. Press 'q' to quit.")

	var b = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		if b[0] == 'q' {
			break
		}
		fmt.Printf("You pressed: %q\n", b[0])
	}

	fmt.Println("Exiting, terminal restored.")
}
