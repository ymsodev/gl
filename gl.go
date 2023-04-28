package gl

type Gl struct {
	env *env
}

func New() *Gl {
	return &Gl{newGlobalEnv()}
}

func (gl *Gl) Init() {

}

func (gl *Gl) Run(code string) glObj {
	tokens := lex(code)
	exprs, err := parse(tokens)
	if err != nil {
		return glErr{err}
	}
	var val glObj = glNil{}
	for _, expr := range exprs {
		val = eval(expr, gl.env)
		if err, ok := val.(glErr); ok {
			return err
		}
	}
	return val
}
