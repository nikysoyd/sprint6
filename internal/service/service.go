package service

import (
	"errors"
	"strings"
)

var morseCode = map[string]string{
	"A": ".-", "B": "-...", "C": "-.-.", "D": "-..", "E": ".",
	"F": "..-.", "G": "--.", "H": "....", "I": "..", "J": ".---",
	"K": "-.-", "L": ".-..", "M": "--", "N": "-.", "O": "---",
	"P": ".--.", "Q": "--.-", "R": ".-.", "S": "...", "T": "-",
	"U": "..-", "V": "...-", "W": ".--", "X": "-..-", "Y": "-.--",
	"Z": "--..", "0": "-----", "1": ".----", "2": "..---", "3": "...--",
	"4": "....-", "5": ".....", "6": "-....", "7": "--...", "8": "---..",
	"9": "----.",
}

var reverseMorse = make(map[string]string)

func init() {
	for k, v := range morseCode {
		reverseMorse[v] = k
	}
}

// Convert определяет тип входной строки (текст или код Морзе) и конвертирует
func Convert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("empty input")
	}

	// Проверяем, является ли строка кодом Морзе
	isMorse := strings.ContainsAny(input, ".-") && !strings.ContainsAny(input, "abcdefghijklmnopqrstuvwxyz0123456789")

	if isMorse {
		return morseToText(input)
	}
	return textToMorse(input)
}

// morseToText конвертирует код Морзе в текст
func morseToText(morse string) (string, error) {
	words := strings.Split(morse, " / ")
	var result strings.Builder

	for i, word := range words {
		chars := strings.Split(word, " ")
		for j, char := range chars {
			if char == "" {
				continue
			}
			if letter, exists := reverseMorse[char]; exists {
				result.WriteString(letter)
			} else {
				return "", errors.New("invalid morse code: " + char)
			}
			if j < len(chars)-1 {
				result.WriteString("")
			}
		}
		if i < len(words)-1 {
			result.WriteString(" ")
		}
	}

	return result.String(), nil
}

// textToMorse конвертирует текст в код Морзе
func textToMorse(text string) (string, error) {
	text = strings.ToUpper(text)
	words := strings.Split(text, " ")
	var result strings.Builder

	for i, word := range words {
		for j, char := range word {
			if code, exists := morseCode[string(char)]; exists {
				result.WriteString(code)
				if j < len(word)-1 {
					result.WriteString(" ")
				}
			} else {
				return "", errors.New("invalid character: " + string(char))
			}
		}
		if i < len(words)-1 {
			result.WriteString(" / ")
		}
	}

	return result.String(), nil
}
