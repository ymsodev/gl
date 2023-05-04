package gl

import (
	"fmt"
)

type env struct {
	outer *env
	data  map[string]GLObject
}

func newEnv(outer *env) *env {
	return &env{outer, make(map[string]GLObject)}
}

func (e *env) set(sym string, val GLObject) {
	e.data[sym] = val
}

func (e *env) get(sym string) (GLObject, error) {
	if res, ok := e.data[sym]; ok {
		return res, nil
	}
	if e.outer != nil {
		return e.outer.get(sym)
	}
	return nil, fmt.Errorf("undefined symbol: %s", sym)
}
