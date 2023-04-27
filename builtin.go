package gl

import "errors"

func add(args ...glObj) glObj {
	res := 0.0
	for _, arg := range args {
		num, ok := arg.(glNum)
		if !ok {
			return glErr{errors.New("operands for + must be numbers")}
		}
		res += num.val
	}
	return glNum{res}
}

func subtract(args ...glObj) glObj {
	res := 0.0
	for _, arg := range args {
		num, ok := arg.(glNum)
		if !ok {
			return glErr{errors.New("operands for - must be numbers")}
		}
		res -= num.val
	}
	return glNum{res}
}

func multiply(args ...glObj) glObj {
	res := 1.0
	for _, arg := range args {
		num, ok := arg.(glNum)
		if !ok {
			return glErr{errors.New("operands for * must be numbers")}
		}
		res *= num.val
	}
	return glNum{res}
}

func divide(args ...glObj) glObj {
	res := 1.0
	for _, arg := range args {
		num, ok := arg.(glNum)
		if !ok {
			return glErr{errors.New("operands for / must be numbers")}
		}
		res /= num.val
	}
	return glNum{res}
}
