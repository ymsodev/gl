package gl

type Gl struct {
	env *env
}

func New() *Gl {
	return &Gl{newEnv(nil)}
}

func (gl *Gl) Init() {

}

func (gl *Gl) Run(code string) ([]any, error) {
	tokens := lex(code)
	exprs, err := parse(tokens)
	if err != nil {
		return nil, err
	}
	vals := make([]any, len(exprs))
	for i, expr := range exprs {
		expr.eval(gl.env)
		vals[i] = expr.value()
	}
	return vals, nil
}
