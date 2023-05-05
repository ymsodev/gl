package gl

import (
	"fmt"
)

type Env struct {
	outer *Env
	data  map[string]GLObject
}

func NewEnv(outer *Env) *Env {
	return &Env{outer, make(map[string]GLObject)}
}

func (e *Env) Set(sym string, val GLObject) {
	e.data[sym] = val
}

func (e *Env) Get(sym string) (GLObject, error) {
	if res, ok := e.data[sym]; ok {
		return res, nil
	}
	if e.outer != nil {
		return e.outer.Get(sym)
	}
	return nil, fmt.Errorf("undefined symbol: %s", sym)
}
