package gl

import (
	"errors"
)

func eval(expr glObj, env *env) glObj {
	switch v := expr.(type) {
	case glSym:
		// TODO: maybe just return the error type?
		res, err := env.get(v.name)
		if err != nil {
			return glErr{err}
		}
		return res
	case glList:
		return evalList(v, env)
	}
	return expr
}

func evalList(l glList, env *env) glObj {
	if len(l.items) == 0 {
		return glNil{}
	}
	car, cdr := l.items[0], l.items[1:]
	if s, ok := car.(glSym); ok {
		switch s.name {
		case "def":
			return def(cdr, env)
		case "let":
			return let(cdr, env)
		}
	}
	for i, item := range l.items {
		val := eval(item, env)
		if err, ok := val.(glErr); ok {
			return err
		}
		l.items[i] = val
	}
	return apply(l, env)
}

func apply(l glList, env *env) glObj {
	f, args := l.items[0], l.items[1:]
	if f, ok := f.(glFn); ok {
		return f.fn(args...)
	}
	return l
}

func def(args []glObj, env *env) glObj {
	if len(args) != 2 {
		return glErr{errors.New("def expects two arguments")}
	}
	sym, ok := args[0].(glSym)
	if !ok {
		return glErr{errors.New("def expects a symbol as the first argument")}
	}
	val := eval(args[1], env)
	if err, ok := val.(glErr); ok {
		return err
	}
	env.set(sym.name, val)
	return val
}

func let(args []glObj, env *env) glObj {
	if len(args) < 2 {
		return glErr{errors.New("let expects at least two arguments")}
	}
	local := newEnv(env)
	params, targ := args[:len(args)-1], args[len(args)-1]
	for _, param := range params {
		tup, ok := param.(glList)
		if !ok || len(tup.items) != 2 {
			return glErr{errors.New("expected a list of two items")}
		}
		sym, ok := tup.items[0].(glSym)
		if !ok {
			return glErr{errors.New("expected a symbol")}
		}
		val := eval(tup.items[1], env)
		if err, ok := val.(glErr); ok {
			return err
		}
		local.set(sym.name, val)
	}
	return eval(targ, local)
}
