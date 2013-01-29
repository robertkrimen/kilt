package kilt

import (
	"unicode"
)

func StringTrimSpace(target string) string {
	// Discard \r? Go already does this for raw string literals.
	end := len(target)

	last := 0
	index := 0
	for index = 0; index < end; index++ {
		chr := rune(target[index])
		if chr == '\n' {
			last = index
		}
		if !unicode.IsSpace(chr) {
			break
		}
	}
	if index >= end {
		return ""
	}
	start := last
	if rune(target[start]) == '\n' {
		// Skip the leading newline
		start++
	}

	last = end - 1
	newline := false
	for index = last; index > start; index-- {
		chr := rune(target[index])
		if chr == '\n' {
			newline = true
		}
		if !unicode.IsSpace(chr) {
			last = index
			break
		}
	}
	stop := last
	if last < end-1 {
		// Include the trailing newline (if any)
		if rune(target[last+1]) == '\n' {
			newline = false // False because slice will already provide a trailing newline
			stop++
		}
	}

	result := target[start : stop+1]
	if newline {
		result += "\n"
	}
	return result
}
