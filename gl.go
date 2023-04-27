package gl

type Gl struct {
	env *env
}

func New() *Gl {
	return &Gl{newGlobalEnv()}
}

func (gl *Gl) Init() {

}

func (gl *Gl) Run(code string) ([]glObj, error) {
	tokens := lex(code)
	exprs, err := parse(tokens)
	if err != nil {
		return nil, err
	}
	vals := make([]glObj, len(exprs))
	for i, expr := range exprs {
		val, err := eval(expr, gl.env)
		if err != nil {
			return nil, err
		}
		vals[i] = val
	}
	return vals, nil
}
