package gl

import (
	"fmt"
)

type env struct {
	outer *env
	data  map[string]glObj
}

func newEnv(outer *env) *env {
	return &env{outer, make(map[string]glObj)}
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
