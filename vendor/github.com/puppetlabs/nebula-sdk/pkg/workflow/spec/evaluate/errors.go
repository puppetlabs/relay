package evaluate

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrUnsupportedLanguage = errors.New("evaluate: unsupported language")
)

type InvalidTypeError struct {
	Type  string
	Cause error
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("could not evaluate a %s type: %+v", e.Type, e.Cause)
}

type FieldNotFoundError struct {
	Name string
}

func (e *FieldNotFoundError) Error() string {
	return fmt.Sprintf("the required field %q could not be found", e.Name)
}

type InvalidEncodingError struct {
	Type  string
	Cause error
}

func (e *InvalidEncodingError) Error() string {
	return fmt.Sprintf("could not evaluate encoding %q: %+v", e.Type, e.Cause)
}

type PathEvaluationError struct {
	Path  string
	Cause error
}

func (e *PathEvaluationError) trace() ([]string, error) {
	var path []string
	for {
		path = append(path, e.Path)

		en, ok := e.Cause.(*PathEvaluationError)
		if !ok {
			return path, e.Cause
		}

		e = en
	}
}

func (e *PathEvaluationError) UnderlyingCause() error {
	_, err := e.trace()
	return err
}

func (e *PathEvaluationError) Error() string {
	path, err := e.trace()
	return fmt.Sprintf("path %q: %+v", strings.Join(path, "."), err)
}

type UnresolvableError struct {
	Causes []error
}

func (e *UnresolvableError) Error() string {
	var causes []string
	for _, err := range e.Causes {
		causes = append(causes, fmt.Sprintf("* %s", err.Error()))
	}

	return fmt.Sprintf("unresolvable:\n%s", strings.Join(causes, "\n"))
}
