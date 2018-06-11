package css

import (
	"fmt"
	"unicode/utf8"
	"strings"
)

const (
	itemError itemType = iota
	itemEOF
	itemProp
	itemValue
)

const eof = 0

type itemType int

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}

	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

type stateFn func(*lexer) stateFn

type lexer struct {
	input string
	start int
	pos   int
	width int
	items chan item
	state stateFn
}

func lex(input string) *lexer {
	l := lexer{
		input: input,
		items: make(chan item, 2),
		state: lexProp,
	}

	return &l
}

func (l *lexer) nextItem() item {
	for {
		select {
		case i := <-l.items:
			return i
		default:
			l.state = l.state(l)
		}
	}
}

func lexProp(l *lexer) stateFn {
	l.acceptRun(";\n \t")
	l.ignore()
	if l.peek() == eof {
		l.emit(itemEOF)
		close(l.items)
		return nil
	}

	if !l.accept("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ-_") {
		return l.errorf("bad property name syntax: it must start with a letter")
	}

	l.acceptRun("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")
	l.emit(itemProp)
	return lexDelimiter
}

func lexDelimiter(l *lexer) stateFn {
	l.acceptRun(" \t")
	l.ignore()
	if l.next() != ':' {
		l.errorf("bad property: value format: expected ':' but not found")
	}
	return lexValue
}

func lexValue(l *lexer) stateFn {
	l.acceptRun(" \t\n")
	l.ignore()

	for {
		switch r := l.peek(); r {
		case '\'':
			if false == l.acceptSingleQuoteString() {
				return nil
			}
		case '"':
			if false == l.acceptDoubleQuoteString() {
				return nil
			}
		case ';':
			l.emit(itemValue)
			l.next()
			l.ignore()
			return lexProp
		case eof:
			if l.pos == l.start {
				return l.errorf("bad value format: expected not empty value")
			}
			l.emit(itemValue)
			l.next()
			l.ignore()
			l.emit(itemEOF)
			return nil
		default:
			l.next()
		}
	}
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() (rune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args),
	}
	return nil
}

func (l *lexer) acceptSingleQuoteString() bool {
	l.next()
	for {
		switch r := l.next(); r {
		case '\'':
			return true
		case '\\':
			l.next()
		case eof:
			l.errorf("unexpected EOF: expected closing single quote", )
			return false
		}
	}
}

func (l *lexer) acceptDoubleQuoteString() bool {
	l.next()
	for {
		switch r := l.next(); r {
		case '"':
			return true
		case '\\':
			l.next()
		case eof:
			l.errorf("unexpected EOF: expected closing double quote", )
			return false
		}
	}
}
