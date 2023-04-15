package gl

import (
	"fmt"
	"strconv"
)

type expr interface {
	eval(env map[string]any) (any, error)
}

type list struct {
	lparen *token
	rparen *token
	items  []expr
}

func (l *list) eval(env map[string]any) (any, error) {
	if len(l.items) == 0 {
		return nil, nil
	}
	vals := make([]any, len(l.items))
	for i, item := range l.items {
		val, err := item.eval(env)
		if err != nil {
			return nil, err
		}
		vals[i] = val
	}
	return nil, nil
}

type atom struct {
	tok *token
}

func (a *atom) eval(env map[string]any) (any, error) {
	switch a.tok.typ {
	case tokNum:
		return strconv.ParseFloat(a.tok.text, 64)
	case tokSym:
		val, ok := env[a.tok.text]
		if !ok {
			return nil, fmt.Errorf("undefined symbol: %s", a.tok.text)
		}
		return val, nil
	case tokStr:
		return a.tok.text[1 : len(a.tok.text)-1], nil
	case tokNil:
		return nil, nil
	case tokTrue:
		return true, nil
	case tokFalse:
		return false, nil
	}
	return nil, nil
}
