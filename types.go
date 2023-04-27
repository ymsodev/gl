package gl

type glObj interface{ glObj() }

type glNil struct{}
type glSym struct{ val string }
type glBool struct{ val bool }
type glNum struct{ val float64 }
type glStr struct{ val string }
type glFn struct{ fn func(...glObj) glObj }
type glList struct{ items []glObj }
type glErr struct{ err error }

func (_ glNil) glObj()  {}
func (_ glSym) glObj()  {}
func (_ glBool) glObj() {}
func (_ glNum) glObj()  {}
func (_ glStr) glObj()  {}
func (_ glFn) glObj()   {}
func (_ glList) glObj() {}
func (_ glErr) glObj()  {}
