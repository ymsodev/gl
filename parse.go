package gl

import (
	"fmt"
)

type expr interface {
	eval(*env) error
	value() any
}

type atom struct {
	tok *token
	val any
}

func newAtom(tok *token) *atom {
	return &atom{tok, nil}
}
func (a *atom) eval(env *env) error { return eval(a, env) }
func (a *atom) value() any          { return a.val }

type list struct {
	lp, rp *token
	items  []expr
	val    any
}

func newList(lp, rp *token, items []expr) *list {
	return &list{lp, rp, items, nil}
}
func (l *list) eval(env *env) error { return eval(l, env) }
func (l *list) value() any          { return l.val }

type parser struct {
	tokens []*token
	exprs  []expr
	curr   int
}

func parse(tokens []*token) ([]expr, error) {
	return newParser(tokens).parse()
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
	case tokNum, tokSym, tokStr:
		return p.atom(), nil
	case tokLParen:
		return p.list()
	default:
		return nil, fmt.Errorf("unexpected token: %s", t.typ)
	}
}

func (p *parser) atom() expr {
	return newAtom(p.next())
}

func (p *parser) list() (expr, error) {
	lp, err := p.consume(tokLParen)
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
	rp, err := p.consume(tokRParen)
	if err != nil {
		return nil, err
	}
	return newList(lp, rp, items), nil
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
