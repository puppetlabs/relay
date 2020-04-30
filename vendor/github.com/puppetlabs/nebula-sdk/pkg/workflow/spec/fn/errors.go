package fn

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrFunctionNotFound          = errors.New("fn: function not found")
	ErrArgNotFound               = errors.New("fn: arg not found")
	ErrPositionalArgsNotAccepted = errors.New("fn: positional arguments cannot be used")
	ErrKeywordArgsNotAccepted    = errors.New("fn: keyword arguments cannot be used")
)

type ArityError struct {
	Wanted   []int
	Variadic bool
	Got      int
}

func (e *ArityError) Error() string {
	wanted := make([]string, len(e.Wanted))
	for i, w := range e.Wanted {
		wanted[i] = strconv.FormatInt(int64(w), 10)
	}

	var variadic string
	if e.Variadic {
		variadic = " or more"
	}

	return fmt.Sprintf("fn: unexpected number of arguments: %d (wanted %s%s)", e.Got, strings.Join(wanted, ", "), variadic)
}

type UnexpectedTypeError struct {
	Wanted []reflect.Type
	Got    reflect.Type
}

func (e *UnexpectedTypeError) Error() string {
	wanted := make([]string, len(e.Wanted))
	for i, w := range e.Wanted {
		wanted[i] = fmt.Sprintf("%s", w)
	}

	return fmt.Sprintf("fn: unexpected type %s (wanted %s)", e.Got, strings.Join(wanted, ", "))
}

type PositionalArgError struct {
	Arg   int
	Cause error
}

func (e *PositionalArgError) Error() string {
	return fmt.Sprintf("fn: arg %d: %+v", e.Arg, e.Cause)
}

type KeywordArgError struct {
	Arg   string
	Cause error
}

func (e *KeywordArgError) Error() string {
	return fmt.Sprintf("fn: arg %q: %+v", e.Arg, e.Cause)
}
