package parser

import (
	"testing"
)

func TestHorizontalRules(t *testing.T) {
	println("Start testing Horizontal Rules")

	parser := getNewParser()
	expected := "<hr/>"

	h1 := "* * *"
	parser.Parse(h1)
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()
	h2 := "***"
	parser.Parse(h2)
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()
	h3 := "- - -"
	parser.Parse(h3)
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()
	h4 := "---"
	parser.Parse(h4)
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()
	h5 := "- -- --- ----"
	parser.Parse(h5)
	PrintError(expected, parser.toHtml(), t)
}

func TestBlockQuotes(t *testing.T) {
	println("Start testing Blockquotes")

	parser := getNewParser()

	q1 := "> 1 Level quotes"
	expected := "<blockquote>1 Level quotes</blockquote>"
	parser.Parse(q1)
	PrintError(expected, parser.toHtml(), t)

	// TODO nested blockquote
}

func TestHeaders(t *testing.T) {
	println("Start testing Headers")

	parser := getNewParser()

	h1 := "# This is an H1"
	parser.Parse(h1)
	expected := "<h1>This is an H1</h1>"
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()

	h2 := "## This is an H2"
	parser.Parse(h2)
	expected = "<h2>This is an H2</h2>"
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()

	h6 := "###### This is an H6"
	parser.Parse(h6)
	expected = "<h6>This is an H6</h6>"
	PrintError(expected, parser.toHtml(), t)
}

func TestLists(t *testing.T) {
	println("Start testing Lists")

	parser := getNewParser()

	list1 := "*   Red"
	parser.Parse(list1)
	expected := "<ul><li>Red</li></ul>"
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()

	list2 := "+   Red"
	parser.Parse(list2)
	expected = "<ul><li>Red</li></ul>"
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()

	list3 := "-   Red"
	parser.Parse(list3)
	expected = "<ul><li>Red</li></ul>"
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()

	list4 := "1.   Red"
	parser.Parse(list4)
	expected = "<ol><li>Red</li></ol>"
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()

	// TODO multi line bullets
	list5 := "2.   Red"
	parser.Parse(list5)
	list5 = "3.   Green"
	parser.Parse(list5)
	list5 = "1.   Blue"
	parser.Parse(list5)
	expected = "<ol><li>Red</li><li>Green</li><li>Blue</li></ol>"
	PrintError(expected, parser.toHtml(), t)
}

func TestCodeBlocks(t *testing.T) {
	println("Start testing about Code Blocks")

	// To produce a code block in Markdown, simply indent every line of the block by at least 4 spaces or 1 tab.
	// For example, given this input:
	parser := getNewParser()

	block1 := "    This is a code block."
	parser.Parse(block1)
	expected := "<pre><code>This is a code block.</code></pre>"
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()

	block2 := "        beep"
	parser.Parse(block2)
	expected = "<pre><code>    beep</code></pre>"
	PrintError(expected, parser.toHtml(), t)

	parser = getNewParser()

	block3 := " beep"
	parser.Parse(block3)
	expected = " beep"
	PrintError(expected, parser.toHtml(), t)
}

func TestLinks(t *testing.T) {
	println("Start testing about Links")

	parser := getNewParser()

	link1 := "[This link](http://example.net/) has no title attribute."
	parser.Parse(link1)
	expected := "<p><a href=\"http://example.net/\">This link</a> has no title attribute.</p>"
	PrintError(expected, parser.toHtml(), t)

	// TODO: Reference-style links
}

func PrintError(expected, actual string, t *testing.T) {
	if actual != expected {
		t.Fatalf("Should be '%v' got '%v'", expected, actual)
	}
}
