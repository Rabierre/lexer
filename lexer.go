package main

import (
  "fmt"
  "strings"
)

type itemType int

const (
  itemError itemType = iota
  itemDot 
  itemEOF 
  itemPipe
  itemLeftMeta
  itemRightMeta
)

type item struct {
  typ itemType
  val string
}

func String(i item) string {
  switch (i.typ) {
  case itemEOF:
    return "EOF"
  case itemError:
    return i.val
  }

  if len(i.val) > 10 {
    return fmt.Sprintf("%.10q", i.val)
  }
  return fmt.Sprintf("%q", i.val)
}

func lex(name, input string) (*lexer, chan item) {
  l := &lexer {
    name: name,
    input: input,
    items: make(chan item),
  }
  
  l.run()
  return l, l.items
}

func main() {
  println(itemDot);
}

type lexer struct {
  name string
  input string
  start int
  pos int
  width int
  items chan item
}

func (l *lexer) run() {
  for state := lexText; state != nil; {
    state = state(l)
  }
  close(l.items)
}

func (l *lexer) emit(t itemType) {
  l.items <- item{t, l.input[l.start:l.pos]}
  l.start = l.pos
}

func lexText(l *lexer) stateFn {
  for {
    if strings.HasPrefix(l.input[l.pos:], leftMeta) {
      if l.pos > l.start {
        l.emit(itemText)
      }
      return lexLeftMeta
    }

    if l.pos > l.start {
      return l.emit(itemText)
    }
    l.emit(itemEOF)
    return nil
  }
}

const leftMeta = "{{"
const rightMeta = "}}"

func lexLeftMeta(l *lexer) stateFn {
  l.pos = len(leftMeta)
  l.emit(itemLeftMeta)
  return lexInsideAction  // Now inside {{ }}
}

func lexInsideAction(l *lexer) stateFn {
  for {
    if strings.HasPrefix(l.input[l.pos:], rightMeta) {
      return lexRightMeta
    }
  }
  
  switch r := l.next(); {
  case r == eof || r == '\n':
    return l.errorf("unclosed action")
  case isSpace(r):
    l.ignore()
  case r == '|':
    l.emit(itemPipe)
  }
}

func (l *lexer) ignore() {
  l.start = l.pos
}

type stateFn struct {
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
  stateFn {
    l.items <- item {
      itemError,
      fmt.Sprintf(format, args...),
    }
  }
  return nil
}
