package main

func IntAbs(a int) int {
	if a < 0 {
		return a * -1
	} else {
		return a
	}
}

func BoolToInt(a bool) int {
	if a {
		return 1
	}

	return 0
}
