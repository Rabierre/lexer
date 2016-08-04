package main

import (
	"fmt"
	"strings"
)

func main() {
	testHeaders()
	testLists()
	testCodeBlocks()
	testLinks()
}

func testHeaders() {
	println("##Starting Test Headers##")

	parser := getNewParser()

	h1 := "# This is an H1"
	parser.Parse(h1)
	println(parser.toHtml() == "<h1>This is an H1</h1>")

	parser = getNewParser()

	h2 := "## This is an H2"
	parser.Parse(h2)
	println(parser.toHtml() == "<h2>This is an H2</h2>")

	parser = getNewParser()

	h6 := "###### This is an H6"
	parser.Parse(h6)
	println(parser.toHtml() == "<h6>This is an H6</h6>")
}

func testLists() {
	println("##Starting Test Lists##")

	parser := getNewParser()

	list1 := "*   Red"
	parser.Parse(list1)
	println(parser.toHtml() == "<ul><li>Red</li></ul>")

	parser = getNewParser()

	list2 := "+   Red"
	parser.Parse(list2)
	println(parser.toHtml() == "<ul><li>Red</li></ul>")

	parser = getNewParser()

	list3 := "-   Red"
	parser.Parse(list3)
	println(parser.toHtml() == "<ul><li>Red</li></ul>")

	parser = getNewParser()

	list4 := "1.   Red"
	parser.Parse(list4)
	println(parser.toHtml() == "<ol><li>Red</li></ol>")

	parser = getNewParser()

	// TODO multi line bullets
	list5 := "1.   Red"
	parser.Parse(list5)
	list5 = "2.   Green"
	parser.Parse(list5)
	println(parser.toHtml() == "<ol><li>Red</li><li>Green</li></ol>")
}

func testCodeBlocks() {
	println("##Starting Test about Code Blocks##")

	// To produce a code block in Markdown, simply indent every line of the block by at least 4 spaces or 1 tab.
	// For example, given this input:
	parser := getNewParser()

	block1 := "    This is a code block."
	parser.Parse(block1)
	println(parser.toHtml() == "<pre><code>This is a code block.</code></pre>")

	parser = getNewParser()

	block2 := "        beep"
	parser.Parse(block2)
	println(parser.toHtml() == "<pre><code>    beep</code></pre>")

	parser = getNewParser()

	block3 := " beep"
	parser.Parse(block3)
	println(parser.toHtml() == " beep")
}

func testLinks() {
	println("##Starting Test about Links##")

	parser := getNewParser()

	link1 := "[This link](http://example.net/) has no title attribute."
	parser.Parse(link1)
	println(parser.toHtml() == "<p><a href=\"http://example.net/\">This link</a> has no title attribute.</p>")
}

func getNewParser() *MarkdownParser {
	return &MarkdownParser{
		contents: []*MarkdownPhrase{},
	}
}

func (parser *MarkdownParser) Parse(input string) *MarkdownItem {
	item := (*MarkdownItem)(nil) // go use typed nil, to initialize with nil need type cast

	if Accept(input, SyntaxHeader) {
		return parser.ParseHeader(input)
	} else if Accept(input, SyntaxListDot) {
		return parser.ParseListDot(input)
	} else if Accept(input, SyntaxListNum) {
		if Accept(string(input[1]), SyntaxDot) {
			return parser.ParseListNumber(input)
		}
	} else if strings.HasPrefix(input, SyntaxCodeBlock1) ||
		strings.HasPrefix(input, SyntaxCodeBlock2) {
		return parser.ParseCodeBlock(input)
	} else if isLink(input) {
		return parser.ParseLink(input)
	} else {
		return parser.ParsePlainText(input)
	}

	return item
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

func Accept(input string, valid string) bool {
	// TODO Accept는 rune단위로 비교하므로 codeblock의 4spaces를 구분할 수 없음
	testChar := rune(input[0])
	if strings.IndexRune(valid, testChar) >= 0 {
		return true
	}
	return false
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
		// slice := lastPrase.items[:len(lastPrase.items)]
		// m.contents[len(m.contents)-1].items = append(slice, item)
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
		if strings.HasPrefix(*input, SpaceSyn) {
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
)

const (
	SyntaxHeader     = "#"
	SpaceSyn         = " "
	SyntaxListDot    = "*+-"
	SyntaxListNum    = "1234567890"
	SyntaxDot        = "."
	SyntaxCodeBlock1 = "    "
	SyntaxCodeBlock2 = "\t"
)
