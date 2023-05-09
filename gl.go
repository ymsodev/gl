package gl

type GL struct {
	ns  Namespace
	env *Environment
}

func New() *GL {
	return &GL{
		ns: Namespace{
			GLSymbol{"+"}: GLFunction{Add},
			GLSymbol{"-"}: GLFunction{Subtract},
			GLSymbol{"*"}: GLFunction{Multiply},
			GLSymbol{"/"}: GLFunction{Divide},
		},
		env: NewEnvironment(nil),
	}
}

func (gl *GL) Init() {
	for sym, fn := range gl.ns {
		gl.env.Set(sym, fn)
	}
}

func (gl *GL) Run(code string) GLObject {
	tokens := Tokenize(code)
	exprs, err := Parse(tokens)
	if err != nil {
		return GLError{err}
	}
	var val GLObject = GLNil{}
	for _, expr := range exprs {
		val = Eval(expr, gl.env)
		if err, ok := val.(GLError); ok {
			return err
		}
	}
	return val
}
