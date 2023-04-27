package gl

import (
	"unicode"
)

type tokType string

const (
	tokLParen tokType = "("
	tokRParen tokType = ")"

	tokSym tokType = "sym"
	tokNum tokType = "num"
	tokStr tokType = "str"

	tokTrue  tokType = "true"
	tokFalse tokType = "false"
	tokNil   tokType = "nil"

	tokEof tokType = "eof"
)

var reserved = map[string]tokType{
	"true":  tokTrue,
	"false": tokFalse,
	"nil":   tokNil,
}

type token struct {
	typ   tokType
	line  int
	start int
	end   int
	text  string
}

func lex(src string) []*token {
	return newLexer(src).lex()
}

type lexer struct {
	src    string
	runes  []rune
	line   int
	start  int
	curr   int
	tokens []*token
}

func newLexer(src string) *lexer {
	return &lexer{
		src:    src,
		runes:  []rune(src),
		line:   0,
		start:  0,
		curr:   0,
		tokens: []*token{},
	}
}

func (l *lexer) lex() []*token {
	for !l.eof() {
		l.scan()
		l.start = l.curr
	}
	l.token(tokEof)
	return l.tokens
}

func (l *lexer) scan() {
	switch r := l.next(); r {
	case '(':
		l.token(tokLParen)
	case ')':
		l.token(tokRParen)
	case ';':
		l.comment()
	case '"':
		l.string()
	case '+', '-':
		if unicode.IsDigit(l.peek()) {
			l.number()
		} else {
			l.token(tokSym)
		}
	default:
		switch {
		case unicode.IsDigit(r):
			l.number()
		case unicode.IsSpace(r):
			if r == '\n' {
				l.line++
			}
		default:
			l.symbol()
		}
	}
}

func (l *lexer) comment() {
	for !l.eof() && l.peek() != '\n' {
		l.next()
	}
}

func (l *lexer) string() {
	for !l.eof() && l.peek() != '"' {
		if r := l.next(); r == '\n' {
			l.line++
		}
	}
	l.next()
	l.token(tokStr)
}

func (l *lexer) number() {
	l.digits()
	if l.peek() == '.' {
		l.next()
		l.digits()
	}
	l.token(tokNum)
}

func (l *lexer) digits() {
	for !l.eof() && unicode.IsDigit(l.peek()) {
		l.next()
	}
}

func (l *lexer) symbol() {
	for !l.eof() && unicode.In(l.peek(), unicode.PrintRanges...) {
		l.next()
	}
	typ := tokSym
	if k, ok := reserved[l.text()]; ok {
		typ = k
	}
	l.token(typ)
}

func (l *lexer) token(typ tokType) {
	l.tokens = append(l.tokens, &token{
		typ:   typ,
		line:  l.line,
		start: l.start,
		end:   l.curr,
		text:  l.text(),
	})
}

func (l *lexer) eof() bool {
	return l.curr >= len(l.runes)
}

func (l *lexer) next() rune {
	defer func() { l.curr++ }()
	return l.peek()
}

func (l *lexer) peek() rune {
	if l.eof() {
		return '\000'
	}
	return l.runes[l.curr]
}

func (l *lexer) text() string {
	return string(l.runes[l.start:l.curr])
}
