package gl

import "strconv"

type evaluator struct {
	env map[string]any
}

func newEvaluator() *evaluator {
	return &evaluator{make(map[string]any)}
}

func (e *evaluator) eval(expr expr) (any, error) {
	switch expr.(type) {
	case *atom:
		return e.evalAtom(expr.(*atom))
	case *list:
		return e.evalList(expr.(*list))
	}
	return nil, nil
}

func (e *evaluator) evalAtom(a *atom) (any, error) {
	switch a.tok.typ {
	case tokNum:
		return strconv.ParseFloat(a.tok.text, 64)
	case tokSym, tokId:
		return e.env[a.tok.text], nil
	case tokStr:
		return a.tok.text[1 : len(a.tok.text)-1], nil
	case tokNil:
		return nil, nil
	case tokTrue:
		return true, nil
	case tokFalse:
		return false, nil
	}
	return nil, nil
}

func (e *evaluator) evalList(l *list) (any, error) {
	vals := []any{}
	for _, item := range l.items {
		val, err := item.eval(e)
	}
}
