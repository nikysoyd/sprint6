package service

import (
	"errors"
	"strings"
	"unicode"

	"github.com/nikysoyd/sprint6/pkg/morse"
)

var (
	ErrEmptyInput = errors.New("empty input")
)

// AutoDetectAndConvert automatically detects whether the input is Morse code or text,
// then converts it to the opposite format.
func AutoDetectAndConvert(input string) (string, error) {
	if len(input) == 0 {
		return "", ErrEmptyInput
	}

	// Trim whitespace from both ends
	input = strings.TrimSpace(input)

	// Check if the input is Morse code (contains only Morse characters: .- or spaces)
	if isMorseCode(input) {
		return morse.ToText(input), nil
	}

	// Otherwise treat as text
	return morse.ToMorse(input), nil
}

// isMorseCode checks if the string appears to be Morse code
func isMorseCode(s string) bool {
	for _, r := range s {
		if !unicode.IsSpace(r) && r != '.' && r != '-' {
			return false
		}
	}
	return true
}
