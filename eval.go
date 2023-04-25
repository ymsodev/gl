package gl

import (
	"fmt"
	"strconv"
)

func eval(expr expr, env *env) error {
	switch v := expr.(type) {
	case *atom:
		return evalAtom(v, env)
	case *list:
		return evalList(v, env)
	}
	return nil
}

func evalAtom(a *atom, env *env) error {
	switch a.tok.typ {
	case tokNum:
		val, err := strconv.ParseFloat(a.tok.text, 64)
		if err != nil {
			return err
		}
		a.val = val
	case tokSym:
		val, err := env.get(a.tok.text)
		if err != nil {
			return err
		}
		a.val = val
	case tokStr:
		a.val = a.tok.text[1 : len(a.tok.text)-1]
	case tokNil:
		a.val = nil
	case tokTrue:
		a.val = true
	case tokFalse:
		a.val = false
	}
	return nil
}

func evalList(l *list, env *env) error {
	if len(l.items) == 0 {
		return nil
	}
	for _, item := range l.items {
		if err := eval(item, env); err != nil {
			return err
		}
	}
	val, err := apply(l, env)
	if err != nil {
		return err
	}
	l.val = val
	return nil
}

func apply(l *list, env *env) (any, error) {
	if car, ok := l.items[0].(*atom); ok {
		cdr := l.items[1:]
		switch car.tok.typ {
		case tokDef:
			return def(cdr, env)
		case tokLet:
			return let(l, env)
		case tokSym:
			if f, ok := car.val.(glFn); ok {
				args := make([]any, len(cdr))
				for i, item := range cdr {
					args[i] = item.value()
				}
				return f.fn(args...), nil
			}
		}
	}
	vals := make([]any, len(l.items))
	for i, item := range l.items {
		vals[i] = item.value()
	}
	return vals, nil
}

func def(args []expr, env *env) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("def expects two arguments")
	}
	s, ok := args[0].(*atom)
	if !ok || s.tok.typ != tokSym {
		return nil, fmt.Errorf("expected a symbol")
	}
	if err := eval(args[1], env); err != nil {
		return nil, err
	}
	val := args[1].value()
	env.set(s.tok.text, val)
	return val, nil
}

func let(l *list, env *env) (any, error) {
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
		if err := eval(tup.items[1], env); err != nil {
			return nil, err
		}
		val := tup.items[1].value()
		local.set(sym.tok.text, val)
	}
	return eval(l.items[len(l.items)-1], local), nil
}
