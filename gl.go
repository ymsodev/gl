package gl

type GL struct {
	env *env
}

func New() *GL {
	return &GL{
		&env{
			nil,
			map[string]GLObject{
				"print": GLFunction{Print},
				"+":     GLFunction{Add},
				"-":     GLFunction{Subtract},
				"*":     GLFunction{Multiply},
				"/":     GLFunction{Divide},
			},
		},
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
