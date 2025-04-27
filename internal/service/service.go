package service

import (
	"errors"
	"strings"
	//"unicode"
)

var morseCode = map[rune]string{
	'А': ".-", 'Б': "-...", 'В': ".--", 'Г': "--.", 'Д': "-..",
	'Е': ".", 'Ж': "...-", 'З': "--..", 'И': "..", 'Й': ".---",
	'К': "-.-", 'Л': ".-..", 'М': "--", 'Н': "-.", 'О': "---",
	'П': ".--.", 'Р': ".-.", 'С': "...", 'Т': "-", 'У': "..-",
	'Ф': "..-.", 'Х': "....", 'Ц': "-.-.", 'Ч': "---.", 'Ш': "----",
	'Щ': "--.-", 'Ъ': "--.--", 'Ы': "-.--", 'Ь': "-..-", 'Э': "..-..",
	'Ю': "..--", 'Я': ".-.-",
	'0': "-----", '1': ".----", '2': "..---", '3': "...--",
	'4': "....-", '5': ".....", '6': "-....", '7': "--...",
	'8': "---..", '9': "----.",
	'.': "......", ',': ".-.-.-", ':': "---...", '?': "..--..",
	';': "-.-.-", '(': "-.--.", ')': "-.--.-", '"': ".-..-.",
	'\'': ".----.", ' ': "/",
}

var reverseMorse = make(map[string]rune)

func init() {
	for k, v := range morseCode {
		reverseMorse[v] = k
	}
}

// Convert определяет тип входной строки и конвертирует
func Convert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", errors.New("пустая строка")
	}

	// Проверяем, является ли строка кодом Морзе
	isMorse := true
	for _, r := range input {
		if !isMorseSymbol(r) && r != ' ' && r != '/' {
			isMorse = false
			break
		}
	}

	if isMorse {
		return morseToText(input)
	}
	return textToMorse(input)
}

func isMorseSymbol(r rune) bool {
	return r == '.' || r == '-'
}

// morseToText конвертирует код Морзе в русский текст
func morseToText(morse string) (string, error) {
	words := strings.Split(morse, " / ")
	var result strings.Builder

	for i, word := range words {
		codes := strings.Split(word, " ")
		for _, code := range codes {
			if code == "" {
				continue
			}
			if char, exists := reverseMorse[code]; exists {
				result.WriteRune(char)
			} else {
				return "", errors.New("некорректный код Морзе: " + code)
			}
		}
		if i < len(words)-1 {
			result.WriteRune(' ')
		}
	}

	return result.String(), nil
}

// textToMorse конвертирует русский текст в код Морзе
func textToMorse(text string) (string, error) {
	text = strings.ToUpper(text)
	var result strings.Builder
	firstChar := true

	for _, char := range text {
		if !firstChar {
			if char == ' ' {
				result.WriteString(" / ")
			} else {
				result.WriteRune(' ')
			}
		}

		if code, exists := morseCode[char]; exists {
			result.WriteString(code)
			firstChar = false
		} else if char != ' ' {
			return "", errors.New("неподдерживаемый символ: " + string(char))
		}
	}

	return result.String(), nil
}
