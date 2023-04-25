package gl

import "fmt"

type env struct {
	data  map[string]any
	outer *env
}

func newEnv(outer *env) *env {
	return &env{make(map[string]any), outer}
}

func (e *env) set(sym string, val any) {
	e.data[sym] = val
}

func (e *env) get(sym string) (any, error) {
	if res, ok := e.data[sym]; ok {
		return res, nil
	}
	if e.outer != nil {
		return e.outer.get(sym)
	}
	return nil, fmt.Errorf("undefined symbol: %s", sym)
}
