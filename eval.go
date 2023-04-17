package gl

import (
	"strconv"
)

func eval(expr expr, env *env) (any, error) {
	switch expr.(type) {
	case *atom:
		return evalAtom(expr.(*atom), env)
	case *list:
		return evalList(expr.(*list), env)
	}
	return nil, nil
}

func evalAtom(a *atom, env *env) (any, error) {
	switch a.tok.typ {
	case tokNum:
		return strconv.ParseFloat(a.tok.text, 64)
	case tokSym:
		return env.get(a.tok.text)
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

func evalList(l *list, env *env) (any, error) {
	if len(l.items) == 0 {
		return nil, nil
	}
	vals := make([]any, len(l.items))
	for i, item := range l.items {
		val, err := eval(item, env)
		if err != nil {
			return nil, err
		}
		vals[i] = val
	}
	return vals, nil
}
