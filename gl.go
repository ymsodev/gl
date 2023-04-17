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
	vals := []any{}
	for _, expr := range exprs {
		expr.eval(gl.env)
		if err := expr.error(); err != nil {
			return vals, err
		}
		vals = append(vals, expr.value())
	}
	return vals, nil
}
