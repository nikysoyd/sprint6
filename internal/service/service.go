package service

import (
	"errors"
	"strings"

	"github.com/nikysoyd/sprint6/pkg/morse"
)

func isMorse(input string) bool {
	isInvalid := func(r rune) bool {
		return !strings.ContainsRune(".-/ ", r)
	}
	return !strings.ContainsFunc(input, isInvalid)
}

func Convert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", morse.ErrNoEncoding{}
	}

	if isMorse(input) {
		result := morse.ToText(input)
		if strings.TrimSpace(result) == "" {
			return "", errors.New("не удалось распознать код Морзе")
		}
		return result, nil
	}

	// Проверка: все ли руны в тексте — буквы, цифры или знаки препинания

	result := morse.ToMorse(input)
	if strings.TrimSpace(result) == "" {
		return "", errors.New("не удалось преобразовать текст в код Морзе")
	}
	return result, nil
}
