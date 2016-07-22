package main

import (
	"./lexer"
	// "bufio"
	// "os"
)

func main() {
	// reader := bufio.NewReader(os.Stdin)
	// text, _ := reader.ReadString('\n')
	// println(text)

	// NUMBER
	result := isNumber("123")
	println(result == true)

	result = isNumber("123a")
	println(result == false)

	result = isNumber("+123")
	println(result == true)

	result = isNumber("-123")
	println(result == true)

	result = isNumber("0x123")
	println(result == true)

	result = isNumber("0X123")
	println(result == true)

	// SPACE
	result = isSpace(" ")
	println(result == true)

	// lexer
	lex := lexer.Lex("1234")
	println(lex.Accept("1234567890") == true)
	println(lex.Accept("abcdefg") == false)

	lex = lexer.Lex("abcd")
	println(lex.Accept("1234567890") == false)
	println(lex.Accept("abcdefg") == true)
}

func isNumber(input string) bool {
	// TODO extract these to accept()
	if rune(input[0]) == '+' || rune(input[0]) == '-' {
		input = input[1:]
	}

	if rune(input[0]) == '0' && (rune(input[1]) == 'x' || rune(input[1]) == 'X') {
		input = input[2:]
	}

	for i := 0; i < len(input); i++ {
		token := rune(input[i])

		if isAlphaNumeric(token) != true {
			return false
		}
	}

	return true
}

func isAlphaNumeric(input rune) bool {
	switch input {
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		return true
	}
	return false
}

func isSpace(input string) bool {
	if input == " " {
		return true
	}
	return false
}
