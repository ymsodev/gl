package gl

import (
	"fmt"
)

type Environment struct {
	outer *Environment
	data  map[GLSymbol]GLObject
}

func NewEnvironment(outer *Environment) *Environment {
	return &Environment{outer, make(map[GLSymbol]GLObject)}
}

func (e *Environment) Set(sym GLSymbol, val GLObject) {
	e.data[sym] = val
}

func (e *Environment) Get(sym GLSymbol) (GLObject, error) {
	if res, ok := e.data[sym]; ok {
		return res, nil
	}
	if e.outer != nil {
		return e.outer.Get(sym)
	}
	return nil, fmt.Errorf("undefined symbol: %s", sym)
}
