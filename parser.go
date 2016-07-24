package main

import (
	"fmt"
	"strings"
)

func main() {
	h1 := "# This is an H1"
	result := Parse(h1)
	println(result.val == "<h1>This is an H1</h1>")
	println(result.typ == itemH1)

	h2 := "## This is an H2"
	result = Parse(h2)
	println(result.val == "<h2>This is an H2</h2>")
	println(result.typ == itemH2)

	h6 := "###### This is an H6"
	result = Parse(h6)
	println(result.val == "<h6>This is an H6</h6>")
	println(result.typ == itemH6)

	list1 := "*   Red"
	result = Parse(list1)
	println(result.val == "<ol><li>Red</li></ol>")
	println(result.typ == itemListDot)

	list2 := "+   Red"
	result = Parse(list2)
	println(result.val == "<ol><li>Red</li></ol>")
	println(result.typ == itemListDot)

	list3 := "-   Red"
	result = Parse(list3)
	println(result.val == "<ol><li>Red</li></ol>")
	println(result.typ == itemListDot)

	list4 := "1.   Red"
	result = Parse(list4)
	println(result.val == "<ul><li>Red</li></ul>")
	println(result.typ == itemListNumber)
}

func Parse(input string) *Item {
	if Accept(input, SyntaxHeader) {
		return ParseHeader(input)
	} else if Accept(input, SyntexListDot) {
		return ParseListDot(input)
	} else if Accept(input, SyntexListNumber) {
		// TODO should check number and dot
	}

	return nil
}

func Accept(input string, valid string) bool {
	testChar := rune(input[0])
	if strings.IndexRune(valid, testChar) >= 0 {
		return true
	}
	return false
}

func ParseListDot(input string) *Item {
	if strings.IndexRune(SpaceSyn, rune(input[1])) >= 0 {
		input = input[1:]

		RemovePrefixSpace(&input)

		trans := fmt.Sprintf("<li>%v</li>", input)
		return &Item{
			val: trans,
			typ: itemListDot,
		}
	}

	return nil
}

func ParseHeader(input string) *Item {
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
	return &Item{
		val: trans,
		typ: itemH1 + MarkdownType(hCnt-1),
	}
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

type Item struct {
	val string
	typ MarkdownType
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
)

const (
	SyntaxHeader  = "#"
	SpaceSyn      = " "
	SyntexListDot = "*+-"
)
