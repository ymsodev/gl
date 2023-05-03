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
			return evalDef(cdr, env)
		case "let":
			return evalLet(cdr, env)
		case "do":
			return evalDo(cdr, env)
		case "if":
			return evalIf(cdr, env)
		case "fn":
			return evalFn(cdr, env)
		}
	}
	vals := make([]glObj, len(l.items))
	for i, item := range l.items {
		val := eval(item, env)
		if err, ok := val.(glErr); ok {
			return err
		}
		vals[i] = val
	}
	return apply(vals, env)
}

func apply(items []glObj, env *env) glObj {
	f, args := items[0], items[1:]
	if f, ok := f.(glFn); ok {
		return f.fn(args...)
	}
	return glList{items}
}

func evalDef(args []glObj, env *env) glObj {
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

func evalLet(args []glObj, env *env) glObj {
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

func evalDo(args []glObj, env *env) glObj {
	var res glObj
	for _, arg := range args {
		res = eval(arg, env)
		if err, ok := res.(glErr); ok {
			return err
		}
	}
	return res
}

func evalIf(args []glObj, env *env) glObj {
	if argc := len(args); argc != 2 && argc != 3 {
		return glErr{errors.New("if expects two or three arguments)")}
	}
	arg0 := eval(args[0], env)
	cond, ok := arg0.(glBool)
	if !ok {
		return glErr{errors.New("expected a bool as a condition")}
	}
	if cond.val {
		return eval(args[1], env)
	}
	if len(args) == 3 {
		return eval(args[2], env)
	}
	return glNil{}
}

func evalFn(args []glObj, env *env) glObj {
	if len(args) != 2 {
		return glErr{errors.New("fn expects two arguments")}
	}
	params, ok := args[0].(glList)
	if !ok {
		return glErr{errors.New("expected a list as parameters")}
	}
	for _, b := range params.items {
		if _, ok := b.(glSym); !ok {
			return glErr{errors.New("expected a list of symbols")}
		}
	}
	body := args[1]
	return glFn{func(args ...glObj) glObj {
		if len(args) != len(params.items) {
			return glErr{errors.New("invalid number of arguments")}
		}
		local := newEnv(env)
		for i, p := range params.items {
			local.set(p.(glSym).name, args[i])
		}
		return eval(body, local)
	}}
}
