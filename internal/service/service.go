package service

import (
	"errors"
	"fmt"
	"strings"
)

// morseCodeMap содержит соответствие символов и кода Морзе
var morseCodeMap = map[string]string{
	"A": ".-", "B": "-...", "C": "-.-.", "D": "-..", "E": ".",
	"F": "..-.", "G": "--.", "H": "....", "I": "..", "J": ".---",
	"K": "-.-", "L": ".-..", "M": "--", "N": "-.", "O": "---",
	"P": ".--.", "Q": "--.-", "R": ".-.", "S": "...", "T": "-",
	"U": "..-", "V": "...-", "W": ".--", "X": "-..-", "Y": "-.--",
	"Z": "--..", "1": ".----", "2": "..---", "3": "...--", "4": "....-",
	"5": ".....", "6": "-....", "7": "--...", "8": "---..", "9": "----.",
	"0": "-----", " ": "/",
}

// reverseMorseCodeMap содержит обратное соответствие кода Морзе и символов
var reverseMorseCodeMap = make(map[string]string)

func init() {
	// Инициализация обратного словаря
	for k, v := range morseCodeMap {
		reverseMorseCodeMap[v] = k
	}
}

// AutoDetectAndConvert автоматически определяет тип ввода и конвертирует
func AutoDetectAndConvert(input string) (string, error) {
	if input == "" {
		return "", errors.New("пустая строка ввода")
	}

	// Определяем тип ввода
	if isMorseCode(input) {
		// Конвертируем код Морзе в текст
		return morseToText(input)
	}
	// Конвертируем текст в код Морзе
	return textToMorse(input)
}

// isMorseCode проверяет, является ли строка кодом Морзе
func isMorseCode(s string) bool {
	for _, r := range s {
		if !(r == '.' || r == '-' || r == ' ' || r == '/') {
			return false
		}
	}
	return true
}

// textToMorse конвертирует текст в код Морзе
func textToMorse(text string) (string, error) {
	var result strings.Builder

	for _, char := range strings.ToUpper(text) {
		strChar := string(char)
		morse, ok := morseCodeMap[strChar]
		if !ok {
			return "", fmt.Errorf("неподдерживаемый символ: %s", strChar)
		}
		result.WriteString(morse)
		result.WriteString(" ")
	}

	return strings.TrimSpace(result.String()), nil
}

// morseToText конвертирует код Морзе в текст
func morseToText(morse string) (string, error) {
	var result strings.Builder

	codes := strings.Split(morse, " ")
	for _, code := range codes {
		if code == "" {
			continue
		}
		char, ok := reverseMorseCodeMap[code]
		if !ok {
			return "", fmt.Errorf("некорректный код Морзе: %s", code)
		}
		result.WriteString(char)
	}

	return result.String(), nil
}
