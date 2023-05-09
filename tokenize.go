package gl

import (
	"unicode"
)

type TokenType byte

const (
	TokLParen TokenType = iota
	TokRParen
	TokSymbol
	TokNumber
	TokString
	TokTrue
	TokFalse
	TokNil
	TokEof
)

func (t TokenType) String() string {
	switch t {
	case TokLParen:
		return "("
	case TokRParen:
		return ")"
	case TokSymbol:
		return "symbol"
	case TokNumber:
		return "number"
	case TokString:
		return "string"
	case TokTrue:
		return "true"
	case TokFalse:
		return "false"
	case TokNil:
		return "nil"
	case TokEof:
		return "EOF"
	default:
		return "<invalid>"
	}
}

var keywords = map[string]TokenType{
	"true":  TokTrue,
	"false": TokFalse,
	"nil":   TokNil,
}

type Token struct {
	Type    TokenType
	Literal string
}

func Tokenize(src string) []Token {
	return newLexer(src).lex()
}

type lexer struct {
	src    string
	runes  []rune
	line   int
	start  int
	curr   int
	tokens []Token
}

func newLexer(src string) *lexer {
	return &lexer{
		src:    src,
		runes:  []rune(src),
		line:   0,
		start:  0,
		curr:   0,
		tokens: []Token{},
	}
}

func (l *lexer) lex() []Token {
	for !l.eof() {
		l.scan()
		l.start = l.curr
	}
	l.token(TokEof)
	return l.tokens
}

func (l *lexer) scan() {
	switch r := l.next(); r {
	case '(':
		l.token(TokLParen)
	case ')':
		l.token(TokRParen)
	case ';':
		l.comment()
	case '"':
		l.string()
	case '+', '-':
		if unicode.IsDigit(l.peek()) {
			l.number()
		} else {
			l.token(TokSymbol)
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
	l.token(TokString)
}

func (l *lexer) number() {
	l.digits()
	if l.peek() == '.' {
		l.next()
		l.digits()
	}
	l.token(TokNumber)
}

func (l *lexer) digits() {
	for !l.eof() && unicode.IsDigit(l.peek()) {
		l.next()
	}
}

func (l *lexer) symbol() {
	for !l.eof() && validSymbolRune(l.peek()) {
		l.next()
	}
	typ := TokSymbol
	if k, ok := keywords[l.text()]; ok {
		typ = k
	}
	l.token(typ)
}

func validSymbolRune(r rune) bool {
	return unicode.In(r, unicode.PrintRanges...) && r != '(' && r != ')'
}

func (l *lexer) token(tokenType TokenType) {
	l.tokens = append(l.tokens, Token{
		Type:    tokenType,
		Literal: l.text(),
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
