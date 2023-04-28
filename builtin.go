package gl

import "errors"

func add(args ...glObj) glObj {
	if len(args) < 2 {
		return glErr{errors.New("+ expects at least two operand")}
	}
	res, ok := args[0].(glNum)
	if !ok {
		return glErr{errors.New("operands for + must be numbers")}
	}
	for _, arg := range args[1:] {
		n, ok := arg.(glNum)
		if !ok {
			return glErr{errors.New("operands for + must be numbers")}
		}
		res.val += n.val
	}
	return res
}

func subtract(args ...glObj) glObj {
	if len(args) < 2 {
		return glErr{errors.New("- expects at least two operand")}
	}
	res, ok := args[0].(glNum)
	if !ok {
		return glErr{errors.New("operands for - must be numbers")}
	}
	for _, arg := range args[1:] {
		num, ok := arg.(glNum)
		if !ok {
			return glErr{errors.New("operands for - must be numbers")}
		}
		res.val -= num.val
	}
	return res
}

func multiply(args ...glObj) glObj {
	if len(args) < 2 {
		return glErr{errors.New("* expects at least two operand")}
	}
	res, ok := args[0].(glNum)
	if !ok {
		return glErr{errors.New("operands for * must be numbers")}
	}
	for _, arg := range args[1:] {
		num, ok := arg.(glNum)
		if !ok {
			return glErr{errors.New("operands for * must be numbers")}
		}
		res.val *= num.val
	}
	return res
}

func divide(args ...glObj) glObj {
	if len(args) < 2 {
		return glErr{errors.New("/ expects at least two operand")}
	}
	res, ok := args[0].(glNum)
	if !ok {
		return glErr{errors.New("operands for / must be numbers")}
	}
	for _, arg := range args[1:] {
		num, ok := arg.(glNum)
		if !ok {
			return glErr{errors.New("operands for / must be numbers")}
		}
		res.val /= num.val
	}
	return res
}
