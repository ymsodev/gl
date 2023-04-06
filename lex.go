package gl

import (
	"fmt"
	"unicode"
)

type tokType byte

const (
	tokLParen tokType = iota
	tokRParen
	tokNum
	tokId
	tokSym
	tokStr
	tokEof
)

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
		case unicode.IsLetter(r):
			l.ident()
		case unicode.IsSpace(r):
			if r == '\n' {
				l.line++
			}
		default:
			l.token(tokSym)
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

func (l *lexer) number() error {
	digits := func() {
		for !l.eof() && unicode.IsDigit(l.peek()) {
			l.next()
		}
	}
	digits()
	if l.peek() == '.' {
		l.next()
		if r := l.peek(); !unicode.IsDigit(r) {
			return fmt.Errorf("expected %#U after .", r)
		}
		digits()
	}
	l.token(tokNum)
	return nil
}

func (l *lexer) ident() {
	valid := func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-'
	}
	for !l.eof() && valid(l.peek()) {
		l.next()
	}
	l.token(tokId)
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
	return l.runes[l.curr]
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
