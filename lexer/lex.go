package lexer

import "strings"

type Item struct {
	typ itemType
	val string
}

type itemType int

const (
	itemError itemType = iota
	itemDot
)

// lexer gets input string and store them in input
// when Run() is called then start lexing input
// generate Item and send them to items channel
type lexer struct {
	name  string
	input string
	start int
	pos   int
	width int
	items chan Item
}

func Lex(_input string) *lexer {
	return &lexer{
		input: _input,
		start: 0,
		pos:   0,
		width: 0,
		items: make(chan Item),
	}
}

// if valid consumes next rune
func (l *lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) next() rune {
	l.pos += l.width

	return rune(l.input[l.pos])
}

func (l *lexer) backup() {
	l.pos -= l.width
}
