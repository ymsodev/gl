package gl

import (
	"errors"
)

func Eval(expr GLObject, env *Environment) GLObject {
	switch v := expr.(type) {
	case GLSymbol:
		res, err := env.Get(v)
		if err != nil {
			return GLError{err}
		}
		return res
	case GLList:
		return evalList(v, env)
	}
	return expr
}

func evalList(l GLList, env *Environment) GLObject {
	if len(l.items) == 0 {
		return GLNil{}
	}
	car, cdr := l.items[0], l.items[1:]
	if s, ok := car.(GLSymbol); ok {
		switch s.name {
		case "def":
			return evalDef(cdr, env)
		case "let":
			return evalLet(cdr, env)
		case "do":
			return evalDo(cdr, env)
		case "if":
			return evalIf(cdr, env)
		case "fn", "\\":
			return evalFn(cdr, env)
		}
	}
	vals := make([]GLObject, len(l.items))
	for i, item := range l.items {
		val := Eval(item, env)
		if err, ok := val.(GLError); ok {
			return err
		}
		vals[i] = val
	}
	return apply(vals, env)
}

func apply(items []GLObject, env *Environment) GLObject {
	f, args := items[0], items[1:]
	if f, ok := f.(GLFunction); ok {
		return f.fn(args...)
	}
	return GLList{items}
}

func evalDef(args []GLObject, env *Environment) GLObject {
	if len(args) != 2 {
		return GLError{errors.New("def expects two arguments")}
	}
	sym, ok := args[0].(GLSymbol)
	if !ok {
		return GLError{errors.New("def expects a symbol as the first argument")}
	}
	val := Eval(args[1], env)
	if err, ok := val.(GLError); ok {
		return err
	}
	env.Set(sym, val)
	return val
}

func evalLet(args []GLObject, env *Environment) GLObject {
	if len(args) < 2 {
		return GLError{errors.New("let expects at least two arguments")}
	}
	local := NewEnvironment(env)
	params, targ := args[:len(args)-1], args[len(args)-1]
	for _, param := range params {
		tup, ok := param.(GLList)
		if !ok || len(tup.items) != 2 {
			return GLError{errors.New("expected a list of two items")}
		}
		sym, ok := tup.items[0].(GLSymbol)
		if !ok {
			return GLError{errors.New("expected a symbol")}
		}
		val := Eval(tup.items[1], env)
		if err, ok := val.(GLError); ok {
			return err
		}
		local.Set(sym, val)
	}
	return Eval(targ, local)
}

func evalDo(args []GLObject, env *Environment) GLObject {
	var res GLObject
	for _, arg := range args {
		res = Eval(arg, env)
		if err, ok := res.(GLError); ok {
			return err
		}
	}
	return res
}

func evalIf(args []GLObject, env *Environment) GLObject {
	if argc := len(args); argc != 2 && argc != 3 {
		return GLError{errors.New("if expects two or three arguments)")}
	}
	arg0 := Eval(args[0], env)
	cond, ok := arg0.(GLBool)
	if !ok {
		return GLError{errors.New("expected a bool as a condition")}
	}
	if cond.val {
		return Eval(args[1], env)
	}
	if len(args) == 3 {
		return Eval(args[2], env)
	}
	return GLNil{}
}

func evalFn(args []GLObject, env *Environment) GLObject {
	if len(args) != 2 {
		return GLError{errors.New("fn expects two arguments")}
	}
	params, ok := args[0].(GLList)
	if !ok {
		return GLError{errors.New("expected a list as parameters")}
	}
	for _, b := range params.items {
		if _, ok := b.(GLSymbol); !ok {
			return GLError{errors.New("expected a list of symbols")}
		}
	}
	body := args[1]
	return GLFunction{func(args ...GLObject) GLObject {
		if len(args) != len(params.items) {
			return GLError{errors.New("invalid number of arguments")}
		}
		local := NewEnvironment(env)
		for i, p := range params.items {
			local.Set(p.(GLSymbol), args[i])
		}
		return Eval(body, local)
	}}
}
