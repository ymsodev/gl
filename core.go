package gl

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errInvalidNumArgs = GLError{errors.New("invalid number of arguments")}
	errInvalidArgType = GLError{errors.New("invalid argument type")}
)

type Namespace map[GLSymbol]GLFunction

func Print(args ...GLObject) GLObject {
	var b strings.Builder
	for i, arg := range args {
		b.WriteString(arg.String())
		if i < len(args)-1 {
			b.WriteRune(' ')
		}
	}
	fmt.Println(b.String())
	return GLNil{}
}

func Add(args ...GLObject) GLObject {
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

func Subtract(args ...GLObject) GLObject {
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

func Multiply(args ...GLObject) GLObject {
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

func Divide(args ...GLObject) GLObject {
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

func List(args ...GLObject) GLObject {
	return GLList{args}
}

func CheckList(args ...GLObject) GLObject {
	if len(args) != 1 {
		return GLError{errors.New("expected 1 argument")}
	}
	_, ok := args[0].(GLList)
	return GLBool{ok}
}

func Len(args ...GLObject) GLObject {
	if len(args) != 1 {
		return GLError{errors.New("expected 1 argument")}
	}
	l, ok := args[0].(GLList)
	if !ok {
		return GLError{errors.New("expected a list as argument")}
	}
	return GLNumber{float64(len(l.items))}
}

func Equal(args ...GLObject) GLObject {
	if len(args) != 2 {
		return GLError{errors.New("expected 2 arguments")}
	}
}
