package gl

import (
	"fmt"
	"strconv"
)

func eval(expr expr, env *env) (any, error) {
	switch v := expr.(type) {
	case *atom:
		return evalAtom(v, env)
	case *list:
		return evalList(v, env)
	default:
		return nil, fmt.Errorf("invalid expr type")
	}
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
	default:
		return nil, fmt.Errorf("invalid token: %s", a.tok.typ)
	}
}

func evalList(l *list, env *env) (any, error) {
	if len(l.items) == 0 {
		return nil, nil
	}
	if car, ok := l.items[0].(*atom); ok {
		switch car.tok.typ {
		case tokDef:
			return evalDef(l, env)
		case tokLet:
			return evalLet(l, env)
		}
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

func evalDef(l *list, env *env) (any, error) {
	if len(l.items) != 3 {
		return nil, fmt.Errorf("def expects two arguments")
	}
	arg1, arg2 := l.items[1], l.items[2]
	s, ok := arg1.(*atom)
	if !ok || s.tok.typ != tokSym {
		return nil, fmt.Errorf("expected a symbol")
	}
	val, err := eval(arg2, env)
	if err != nil {
		return nil, err
	}
	env.set(s.tok.text, val)
	return val, nil
}

func evalLet(l *list, env *env) (any, error) {
	if len(l.items) < 3 {
		return nil, fmt.Errorf("let expects at least two arguments")
	}
	local := newEnv(env)
	for _, item := range l.items[1 : len(l.items)-1] {
		tup, ok := item.(*list)
		if !ok || len(tup.items) != 2 {
			return nil, fmt.Errorf("expected a list of two items")
		}
		sym, ok := tup.items[0].(*atom)
		if !ok || sym.tok.typ != tokSym {
			return nil, fmt.Errorf("expected a symbol")
		}
		val, err := eval(tup.items[1], env)
		if err != nil {
			return nil, err
		}
		local.set(sym.tok.text, val)
	}
	return eval(l.items[len(l.items)-1], local)
}
