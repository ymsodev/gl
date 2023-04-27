package gl

import (
	"fmt"
	"strconv"
)

func parse(tokens []*token) ([]glObj, error) {
	return newParser(tokens).parse()
}

type parser struct {
	tokens []*token
	exprs  []glObj
	curr   int
}

func newParser(tokens []*token) *parser {
	return &parser{
		tokens: tokens,
		exprs:  []glObj{},
		curr:   0,
	}
}

func (p *parser) parse() ([]glObj, error) {
	for !p.eof() {
		expr, err := p.expr()
		if err != nil {
			return p.exprs, err
		}
		p.exprs = append(p.exprs, expr)
	}
	return p.exprs, nil
}

func (p *parser) expr() (glObj, error) {
	switch t := p.peek(); t.typ {
	case tokLParen:
		return p.list()
	default:
		return p.atom()
	}
}

func (p *parser) list() (glObj, error) {
	if _, err := p.consume(tokLParen); err != nil {
		return nil, err
	}
	items := []glObj{}
	for !p.check(tokRParen) {
		item, err := p.expr()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if _, err := p.consume(tokRParen); err != nil {
		return nil, err
	}
	return glList{items}, nil
}

func (p *parser) atom() (glObj, error) {
	switch t := p.next(); t.typ {
	case tokSym:
		return glSym{t.text}, nil
	case tokTrue:
		return glBool{true}, nil
	case tokFalse:
		return glBool{false}, nil
	case tokNum:
		val, err := strconv.ParseFloat(t.text, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %s", t.text)
		}
		return glNum{val}, nil
	case tokStr:
		return glStr{t.text[1 : len(t.text)-1]}, nil
	case tokNil:
		return glNil{}, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", t.typ)
	}
}

func (p *parser) eof() bool {
	return p.peek().typ == tokEof
}

func (p *parser) peek() *token {
	return p.tokens[p.curr]
}

func (p *parser) check(typ tokType) bool {
	return p.peek().typ == typ
}

func (p *parser) next() *token {
	defer func() { p.curr++ }()
	return p.peek()
}

func (p *parser) consume(typ tokType) (*token, error) {
	if !p.check(typ) {
		return nil, fmt.Errorf("expected %s", typ)
	}
	return p.next(), nil
}
