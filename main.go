package main

import (
	"fmt"
	"os"
)

func main() {
	term := NewTerminal()
	term.EnableRawMode()
	term.EnableAlternateBuffer()
	term.ClearAlternateBuffer()
	defer term.Restore()

	fmt.Println("Raw mode enabled. Press 'q' to quit.")

	var b = make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)

		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		if b[0] == 'q' {
			break
		}
		fmt.Printf("You pressed: %q\n", b[0])
	}

	fmt.Println("Exiting, terminal restored.")
}
