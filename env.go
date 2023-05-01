package gl

import (
	"fmt"
)

type env struct {
	outer *env
	binds map[string]glObj
	data  map[string]glObj
}

func newEnv(outer *env, binds map[string]glObj) *env {
	return &env{outer, binds, make(map[string]glObj)}
}

func newGlobalEnv() *env {
	return &env{
		nil,
		map[string]glObj{},
		map[string]glObj{
			"+": glFn{add},
			"-": glFn{subtract},
			"*": glFn{multiply},
			"/": glFn{divide},
		},
	}
}

func (e *env) set(sym string, val glObj) {
	e.data[sym] = val
}

func (e *env) get(sym string) (glObj, error) {
	if res, ok := e.data[sym]; ok {
		return res, nil
	}
	if e.outer != nil {
		return e.outer.get(sym)
	}
	return nil, fmt.Errorf("undefined symbol: %s", sym)
}
