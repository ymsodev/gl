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

func (e *env) find(sym string) *env {
	if _, ok := e.data[sym]; ok {
		return e
	}
	if e.outer == nil {
		return nil
	}
	return e.outer.find(sym)
}

func (e *env) get(sym string) (any, error) {
	if env := e.find(sym); env != nil {
		return env.data[sym], nil
	}
	return nil, fmt.Errorf("undefined symbol: %s", sym)
}
