package gl

import (
	"fmt"
)

func parse(tokens []*token) ([]expr, error) {
	return newParser(tokens).parse()
}

type parser struct {
	tokens []*token
	exprs  []expr
	curr   int
}

func newParser(tokens []*token) *parser {
	return &parser{
		tokens: tokens,
		exprs:  []expr{},
		curr:   0,
	}
}

func (p *parser) parse() ([]expr, error) {
	for !p.eof() {
		expr, err := p.expr()
		if err != nil {
			return p.exprs, err
		}
		p.exprs = append(p.exprs, expr)
	}
	return p.exprs, nil
}

func (p *parser) expr() (expr, error) {
	switch t := p.peek(); t.typ {
	case tokLParen:
		return p.list()
	case tokNum, tokId, tokSym, tokStr:
		return p.atom(), nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", t.typ)
	}
}

func (p *parser) list() (expr, error) {
	lparen, err := p.consume(tokLParen)
	if err != nil {
		return nil, err
	}
	items := []expr{}
	for !p.check(tokRParen) {
		item, err := p.expr()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	rparen, err := p.consume(tokRParen)
	if err != nil {
		return nil, err
	}
	return &list{lparen, rparen, items}, nil
}

func (p *parser) atom() expr {
	return &atom{p.next()}
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
