package typeutil

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type InvalidVersionKindError struct {
	ExpectedVersion, ExpectedKind string
	GotVersion, GotKind           string
}

func (e *InvalidVersionKindError) Error() string {
	m := fmt.Sprintf("typeutil: expected version %q with kind %q", e.ExpectedVersion, e.ExpectedKind)
	if e.GotVersion != "" || e.GotKind != "" {
		m += fmt.Sprintf(", but got version %q with kind %q", e.GotVersion, e.GotKind)
	}

	return m
}

type FieldValidationError struct {
	Context     string
	Field       string
	Description string
	Type        string
}

func (e *FieldValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Description)
}

type ValidationError struct {
	FieldErrors []*FieldValidationError
}

func (e *ValidationError) Error() string {
	if len(e.FieldErrors) == 0 {
		return "typeutil: general validation error"
	}

	fes := make([]string, len(e.FieldErrors))
	for i, fe := range e.FieldErrors {
		fes[i] = fe.Error()
	}

	return fmt.Sprintf("typeutil: validation error:\n* %s", strings.Join(fes, "\n* "))
}

// The following jsonschema field error types should be excluded from the constructed
// ValidationError because they refer to internal abstractions used in the schema document
// such as anyOf, oneOf, allOf, if / else rather than the end-user document being validated.
// Each will be accompanied by an additional error more specifically referencing the issue
// in the document being validated
var excludedResultErrorTypes = map[string]bool{
	"number_any_of":  true,
	"number_one_of":  true,
	"number_all_of":  true,
	"number_not":     true,
	"condition_then": true,
	"condition_else": true,
}

func ValidationErrorFromResult(result *gojsonschema.Result) error {
	if result.Valid() {
		return nil
	}

	errs := result.Errors()

	var fes []*FieldValidationError

	for _, err := range errs {
		if !excludedResultErrorTypes[err.Type()] {
			fes = append(fes, &FieldValidationError{
				Context:     err.Context().String(),
				Field:       err.Field(),
				Description: err.Description(),
				Type:        err.Type(),
			})
		}
	}

	return &ValidationError{
		FieldErrors: fes,
	}
}
