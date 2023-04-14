package gl

type Gl struct {
	eval *evaluator
}

func New() *Gl {
	return &Gl{newEvaluator()}
}

func (gl *Gl) Init() {

}

func (gl *Gl) Run(code string) error {
	tokens := lex(code)
	exprs, err := parse(tokens)
	if err != nil {
		return err
	}
	for _, expr := range exprs {
		if err := expr.eval(gl.eval); err != nil {
			return err
		}
	}
	return nil
}
