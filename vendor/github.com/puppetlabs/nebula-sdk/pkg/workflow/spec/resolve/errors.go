package resolve

import (
	"fmt"
)

type DataQueryError struct {
	Query string
}

func (e *DataQueryError) Error() string {
	return fmt.Sprintf("resolve: data query %q could not be processed", e.Query)
}

type DataNotFoundError struct {
	Query string
}

func (e *DataNotFoundError) Error() string {
	return fmt.Sprintf("resolve: data for query %q could not be found", e.Query)
}

type SecretNotFoundError struct {
	Name string
}

func (e *SecretNotFoundError) Error() string {
	return fmt.Sprintf("resolve: secret %q could not be found", e.Name)
}

type ConnectionNotFoundError struct {
	Type string
	Name string
}

func (e *ConnectionNotFoundError) Error() string {
	return fmt.Sprintf("resolve: connection type %q with name %q could not be found", e.Type, e.Name)
}

type OutputNotFoundError struct {
	From string
	Name string
}

func (e *OutputNotFoundError) Error() string {
	return fmt.Sprintf("resolve: output %q of step %q could not be found", e.Name, e.From)
}

type ParameterNotFoundError struct {
	Name string
}

func (e *ParameterNotFoundError) Error() string {
	return fmt.Sprintf("resolve: parameter %q could not be found", e.Name)
}

type AnswerNotFoundError struct {
	AskRef string
	Name   string
}

func (e *AnswerNotFoundError) Error() string {
	return fmt.Sprintf("resolve: answer %q of ask %q could not be found", e.Name, e.AskRef)
}

type FunctionResolutionError struct {
	Name  string
	Cause error
}

func (e *FunctionResolutionError) Error() string {
	return fmt.Sprintf("resolve: function %q could not be invoked: %+v", e.Name, e.Cause)
}
