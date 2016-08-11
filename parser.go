package parser

import (
	"fmt"
	"strings"
)

func getNewParser() *MarkdownParser {
	return &MarkdownParser{
		contents: []*MarkdownPhrase{},
	}
}

func (parser *MarkdownParser) Parse(input string) {
	if Accept(input, SyntaxHeader) {
		parser.ParseHeader(input)
	} else if isHorizontalRule(input) {
		parser.ParseHorizontalRule(input)
	} else if isList(input) {
		parser.ParseList(input)
	} else if strings.HasPrefix(input, SyntaxCodeBlock1) ||
		strings.HasPrefix(input, SyntaxCodeBlock2) {
		parser.ParseCodeBlock(input)
	} else if isLink(input) {
		parser.ParseLink(input)
	} else if Accept(input, SyntaxBlockquote) {
		parser.ParseBlockquote(input)
	} else {
		parser.ParsePlainText(input)
	}
}

func isList(input string) bool {
	return Accept(input, SyntaxListDot) || (Accept(input, SyntaxListNum) && Accept(string(input[1]), SyntaxDot))
}

func isLink(input string) bool {
	// [This link](http://example.net/) has no title attribute.
	openBracket := getIndexFromString("[", input, 0)
	closeBracket := getIndexFromString("]", input, openBracket)
	openParenthese := getIndexFromString("(", input, closeBracket)
	closeParenthese := getIndexFromString(")", input, openParenthese)

	if openBracket < closeBracket && openParenthese < closeParenthese {
		return true
	}
	return false
}

func isHorizontalRule(input string) bool {
	for _, token := range input {
		if string(token) == SyntaxSpace {
			continue
		}

		if string(token) != "*" && string(token) != "-" {
			return false
		}
	}
	return true
}

func Accept(input string, valid string) bool {
	// TODO Accept는 rune단위로 비교하므로 codeblock의 4spaces를 구분할 수 없음
	testChar := rune(input[0])
	if strings.IndexRune(valid, testChar) >= 0 {
		return true
	}
	return false
}

func (m *MarkdownParser) ParseBlockquote(input string) {
	if strings.HasPrefix(input, SyntaxBlockquote) {
		input = input[1:]
	}

	RemovePrefixSpace(&input)

	item := &MarkdownItem{
		val: input,
		typ: itemBlockquote,
	}

	m.addToParsedList(item)
}

func (m *MarkdownParser) ParseList(input string) {
	if Accept(input, SyntaxListDot) {
		m.ParseListDot(input)
	} else {
		m.ParseListNumber(input)
	}
}

func (m *MarkdownParser) ParseLink(input string) *MarkdownItem {
	// TODO invalid can be parsed
	openBracket := getIndexFromString("[", input, 0)
	closeBracket := getIndexFromString("]", input, openBracket)
	openParenthese := getIndexFromString("(", input, closeBracket)
	closeParenthese := getIndexFromString(")", input, openParenthese)

	link := input[openParenthese+1 : closeParenthese]
	linkedText := input[openBracket+1 : closeBracket]
	left := input[closeParenthese+1:]
	trans := fmt.Sprintf("<p><a href=\"%v\">%v</a>%v</p>", link, linkedText, left)

	item := &MarkdownItem{
		val: trans,
		typ: itemLink,
	}

	m.addToParsedList(item)

	return item
}

func getIndexFromString(valid string, input string, index int) int {
	for i := index; i >= 0 && i < len(input); i++ {
		if strings.IndexRune(valid, rune(input[i])) >= 0 {
			return i
		}
	}
	return -1
}

func (m *MarkdownParser) ParseHorizontalRule(input string) {
	m.addToParsedList(&MarkdownItem{
		val: "<hr/>",
		typ: itemHorizontalRule,
	})
}

func (m *MarkdownParser) ParsePlainText(input string) *MarkdownItem {
	item := &MarkdownItem{
		val: input,
		typ: itemPlainText,
	}

	m.addToParsedList(item)

	return item
}

// TODO: parse as block
func (m *MarkdownParser) ParseCodeBlock(input string) *MarkdownItem {
	if strings.HasPrefix(input, SyntaxCodeBlock1) {
		input = input[len(SyntaxCodeBlock1):]
	} else if strings.HasPrefix(input, SyntaxCodeBlock2) {
		input = input[len(SyntaxCodeBlock2):]
	}

	trans := fmt.Sprintf("<code>%v</code>", input)
	item := &MarkdownItem{
		val: trans,
		typ: itemCodeBlock,
	}

	m.addToParsedList(item)

	return item
}

func (m *MarkdownParser) ParseListNumber(input string) *MarkdownItem {
	if strings.IndexRune(SyntaxListNum, rune(input[0])) >= 0 &&
		strings.IndexRune(SyntaxDot, rune(input[1])) == 0 {
		input = input[2:]
	}

	RemovePrefixSpace(&input)

	trans := fmt.Sprintf("<li>%v</li>", input)
	item := &MarkdownItem{
		val: trans,
		typ: itemListNumber,
	}

	m.addToParsedList(item)

	return item
}

func (m *MarkdownParser) addToParsedList(item *MarkdownItem) {
	if len(m.contents) == 0 {
		mp := &MarkdownPhrase{
			typ:   item.typ,
			items: []*MarkdownItem{item},
		}
		m.AddNewPhrase(mp)
	} else if lastPrase := m.contents[len(m.contents)-1]; lastPrase.typ == item.typ {
		lastPrase.AddNewItem(item)
	} else {
		mp := &MarkdownPhrase{
			typ:   item.typ,
			items: []*MarkdownItem{item},
		}
		m.AddNewPhrase(mp)
	}
}

func (m *MarkdownPhrase) AddNewItem(item *MarkdownItem) {
	// m.contents = append(m.contents[:len(m.contents)], mp)
	slice := m.items[:len(m.items)]
	m.items = append(slice, item)
}

func (m *MarkdownParser) AddNewPhrase(mp *MarkdownPhrase) {
	m.contents = append(m.contents[:len(m.contents)], mp)
}

// TODO: parse as block
func (m *MarkdownParser) ParseListDot(input string) *MarkdownItem {
	if strings.IndexRune(SyntaxListDot, rune(input[0])) >= 0 {
		input = input[1:]
	}

	RemovePrefixSpace(&input)

	trans := fmt.Sprintf("<li>%v</li>", input)
	item := &MarkdownItem{
		val: trans,
		typ: itemListDot,
	}

	m.addToParsedList(item)

	return item
}

func (m *MarkdownParser) ParseHeader(input string) *MarkdownItem {
	hCnt := 0
	for {
		if strings.HasPrefix(input, SyntaxHeader) {
			hCnt++
			input = input[1:]
		} else {
			break
		}
	}

	RemovePrefixSpace(&input)

	trans := fmt.Sprintf("<h%v>%v</h%v>", hCnt, input, hCnt)
	item := &MarkdownItem{
		val: trans,
		typ: itemH1 + MarkdownType(hCnt-1),
	}

	m.addToParsedList(item)

	return item
}

func RemovePrefixSpace(input *string) {
	for {
		if strings.HasPrefix(*input, SyntaxSpace) {
			// slice returns ptr of string
			*input = (*input)[1:]
		} else {
			break
		}
	}
}

type MarkdownParser struct {
	contents []*MarkdownPhrase
}

func (m *MarkdownParser) toHtml() string {
	html := ""
	for _, con := range m.contents {
		html += con.toHtml()
	}

	return html
}

type MarkdownPhrase struct {
	typ   MarkdownType
	items []*MarkdownItem
}

func (m *MarkdownPhrase) toHtml() string {
	html := m.getPrefix()
	for _, item := range m.items {
		html += item.toHtml()
	}

	html += m.getSufix()

	return html
}

func (m *MarkdownPhrase) getPrefix() string {
	if m.typ == itemListDot {
		return "<ul>"
	} else if m.typ == itemListNumber {
		return "<ol>"
	} else if m.typ == itemCodeBlock {
		return "<pre>"
	} else if m.typ == itemBlockquote {
		return "<blockquote>"
	}
	return ""
}

func (m *MarkdownPhrase) getSufix() string {
	if m.typ == itemListDot {
		return "</ul>"
	} else if m.typ == itemListNumber {
		return "</ol>"
	} else if m.typ == itemCodeBlock {
		return "</pre>"
	} else if m.typ == itemBlockquote {
		return "</blockquote>"
	}
	return ""
}

type MarkdownItem struct {
	val string
	typ MarkdownType
}

func (m *MarkdownItem) toHtml() string {
	return m.val
}

type MarkdownType int

const (
	_ MarkdownType = iota
	itemH1
	itemH2
	itemH3
	itemH4
	itemH5
	itemH6
	itemListDot
	itemListNumber
	itemCodeBlock
	itemLink
	itemPlainText
	itemBlockquote
	itemHorizontalRule
)

const (
	SyntaxHeader     = "#"
	SyntaxSpace      = " "
	SyntaxListDot    = "*+-"
	SyntaxListNum    = "1234567890"
	SyntaxDot        = "."
	SyntaxCodeBlock1 = "    "
	SyntaxCodeBlock2 = "\t"
	SyntaxBlockquote = ">"
)
