package gl

import (
	"fmt"
	"strconv"
)

func Parse(tokens []*Token) ([]GLObject, error) {
	return newParser(tokens).parse()
}

type parser struct {
	tokens []*Token
	exprs  []GLObject
	curr   int
}

func newParser(tokens []*Token) *parser {
	return &parser{
		tokens: tokens,
		exprs:  []GLObject{},
		curr:   0,
	}
}

func (p *parser) parse() ([]GLObject, error) {
	for !p.eof() {
		expr, err := p.expr()
		if err != nil {
			return p.exprs, err
		}
		p.exprs = append(p.exprs, expr)
	}
	return p.exprs, nil
}

func (p *parser) expr() (GLObject, error) {
	switch t := p.peek(); t.typ {
	case TokLeftParen:
		return p.list()
	default:
		return p.atom()
	}
}

func (p *parser) list() (GLObject, error) {
	if _, err := p.consume(TokLeftParen); err != nil {
		return nil, err
	}
	items := []GLObject{}
	for !p.check(TokRightParen) {
		item, err := p.expr()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if _, err := p.consume(TokRightParen); err != nil {
		return nil, err
	}
	return GLList{items}, nil
}

func (p *parser) atom() (GLObject, error) {
	switch t := p.next(); t.typ {
	case TokSymbol:
		return GLSymbol{t.text}, nil
	case TokTrue:
		return GLBool{true}, nil
	case TokFalse:
		return GLBool{false}, nil
	case TokNumber:
		val, err := strconv.ParseFloat(t.text, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %s", t.text)
		}
		return GLNumber{val}, nil
	case TokString:
		return GLString{t.text[1 : len(t.text)-1]}, nil
	case TokNil:
		return GLNil{}, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", t.typ)
	}
}

func (p *parser) eof() bool {
	return p.peek().typ == TokEof
}

func (p *parser) peek() *Token {
	return p.tokens[p.curr]
}

func (p *parser) check(typ TokenType) bool {
	return p.peek().typ == typ
}

func (p *parser) next() *Token {
	defer func() { p.curr++ }()
	return p.peek()
}

func (p *parser) consume(typ TokenType) (*Token, error) {
	if !p.check(typ) {
		return nil, fmt.Errorf("expected %s", typ)
	}
	return p.next(), nil
}
