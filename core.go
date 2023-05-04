package gl

import "errors"

var (
	errInvalidNumArgs = GLError{errors.New("invalid number of arguments")}
	errInvalidArgType = GLError{errors.New("invalid argument type")}
)

func add(args ...GLObject) GLObject {
	if len(args) < 2 {
		return errInvalidNumArgs
	}
	res, ok := args[0].(GLNumber)
	if !ok {
		return errInvalidArgType
	}
	for _, arg := range args[1:] {
		n, ok := arg.(GLNumber)
		if !ok {
			return errInvalidArgType
		}
		res.val += n.val
	}
	return res
}

func subtract(args ...GLObject) GLObject {
	if len(args) < 2 {
		return errInvalidNumArgs
	}
	res, ok := args[0].(GLNumber)
	if !ok {
		return errInvalidArgType
	}
	for _, arg := range args[1:] {
		num, ok := arg.(GLNumber)
		if !ok {
			return errInvalidArgType
		}
		res.val -= num.val
	}
	return res
}

func multiply(args ...GLObject) GLObject {
	if len(args) < 2 {
		return errInvalidNumArgs
	}
	res, ok := args[0].(GLNumber)
	if !ok {
		return errInvalidArgType
	}
	for _, arg := range args[1:] {
		num, ok := arg.(GLNumber)
		if !ok {
			return errInvalidArgType
		}
		res.val *= num.val
	}
	return res
}

func divide(args ...GLObject) GLObject {
	if len(args) < 2 {
		return errInvalidNumArgs
	}
	res, ok := args[0].(GLNumber)
	if !ok {
		return errInvalidArgType
	}
	for _, arg := range args[1:] {
		num, ok := arg.(GLNumber)
		if !ok {
			return errInvalidArgType
		}
		res.val /= num.val
	}
	return res
}
