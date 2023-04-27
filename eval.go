package gl

import (
	"fmt"
)

func eval(expr glObj, env *env) (glObj, error) {
	switch v := expr.(type) {
	case glSym:
		return env.get(v.val)
	case glList:
		return evalList(v, env)
	}
	return expr, nil
}

func evalList(l glList, env *env) (glObj, error) {
	if len(l.items) == 0 {
		return glNil{}, nil
	}
	if car, ok := l.items[0].(*atom); ok {
		cdr := l.items[1:]
		switch car.tok.typ {
		case tokDef:
			val, err := def(cdr, env)
			if err != nil {
				return err
			}
			l.val = val
			return nil

		}
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
			return let(cdr, env)
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
	val := args[1].value()
	env.set(s.tok.text, val)
	return val, nil
}

func let(args []expr, env *env) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("let expects at least two arguments")
	}
	local := newEnv(env)
	for _, arg := range args[:len(args)-1] {
		tup, ok := arg.(*list)
		if !ok || len(tup.items) != 2 {
			return nil, fmt.Errorf("expected a list of two items")
		}
		sym, ok := tup.items[0].(*atom)
		if !ok || sym.tok.typ != tokSym {
			return nil, fmt.Errorf("expected a symbol")
		}
		val := tup.items[1].value()
		local.set(sym.tok.text, val)
	}
	return eval(args[len(args)-1], local), nil
}
