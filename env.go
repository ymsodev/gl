package gl

import (
	"fmt"
)

type env struct {
	data  map[string]glObj
	outer *env
}

func newEnv(outer *env) *env {
	return &env{make(map[string]glObj), outer}
}

func newGlobalEnv() *env {
	return &env{
		map[string]glObj{
			"+": glFn{add},
			"-": glFn{subtract},
			"*": glFn{multiply},
			"/": glFn{divide},
		},
		nil,
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
