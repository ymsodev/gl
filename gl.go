package gl

import "fmt"

type Gl struct {
	env *env
}

func New() *Gl {
	return &Gl{
		&env{
			nil,
			map[string]glObj{
				"+": glFn{add},
				"-": glFn{subtract},
				"*": glFn{multiply},
				"/": glFn{divide},
			},
		},
	}
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

func Print(obj glObj) {
	switch v := obj.(type) {
	case glNil:
		fmt.Println("nil")
	case glSym:
		fmt.Println(v.name)
	case glBool:
		fmt.Println(v.val)
	case glNum:
		fmt.Println(v.val)
	case glStr:
		fmt.Println(v.val)
	case glFn:
		fmt.Println("function")
	case glList:
		fmt.Println("(")
		for _, item := range v.items {
			fmt.Print("\t")
			Print(item)
		}
		fmt.Println(")")
	case glErr:
		fmt.Println(v.err)
	}
}
