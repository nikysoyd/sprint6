package service

import (
	"errors"
	"strings"

	"github.com/nikysoyd/sprint6/pkg/morse"
)

func isMorse(input string) bool {
	isletter := func(r rune) bool {
		return !strings.ContainsRune(".-/ ", r)
	}
	return !strings.ContainsFunc(input, isletter)
}

func Convert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", morse.ErrNoEncoding{}
	}

	if isMorse(input) {
		result := morse.ToText(input)
		if strings.TrimSpace(result) == "" {
			return "", morse.ErrNoEncoding{}
		}
		return result, nil
	}

	result := morse.ToMorse(input)
	if strings.TrimSpace(result) == "" {
		return "", errors.New("error converting code to Morse")
	}
	return result, nil
}
