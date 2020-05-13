package generator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrImageNotFound = errors.New("generator: image definition not found")
)

type ImageDependencySelfReferenceError struct {
	ImageName string
}

func (e *ImageDependencySelfReferenceError) Error() string {
	return fmt.Sprintf("generator: image %q cannot depend on itself", e.ImageName)
}

type ImageDependencyMissingError struct {
	ImageName string
	Want      string
}

func (e *ImageDependencyMissingError) Error() string {
	return fmt.Sprintf("generator: image %q depends on %q, but the dependent image does not exist in this manifest", e.ImageName, e.Want)
}

type ImageDependencyCyclesError struct {
	Cycles [][]string
}

func (e *ImageDependencyCyclesError) Error() string {
	cycles := make([]string, len(e.Cycles))
	for i, cycle := range e.Cycles {
		cycles[i] = strings.Join(cycle, " -> ") + " -> " + cycle[0]
	}

	if len(cycles) == 1 {
		return fmt.Sprintf("generator: cycles detected in images: %s", cycles[0])
	}

	return fmt.Sprintf("generator: cycles detected in images:\n  %s", strings.Join(cycles, "\n  "))
}

type ImageTemplateParseError struct {
	ImageName string
	Cause     error
}

func (e *ImageTemplateParseError) Error() string {
	return fmt.Sprintf("generator: error parsing template for %q: %+v", e.ImageName, e.Cause)
}

type ImageTemplateExecutionError struct {
	ImageName string
	Cause     error
}

func (e *ImageTemplateExecutionError) Error() string {
	return fmt.Sprintf("generator: error executing template for %q: %+v", e.ImageName, e.Cause)
}

type ImageTemplateFormatError struct {
	ImageName string
	Content   string
	Cause     error
}

func (e *ImageTemplateFormatError) Error() string {
	return fmt.Sprintf("generator: template for %q did not produce a valid Dockerfile: %+v", e.ImageName, e.Cause)
}

type ImageTemplateNoStagesError struct {
	ImageName string
	Content   string
}

func (e *ImageTemplateNoStagesError) Error() string {
	return fmt.Sprintf("generator: template for %q did not produce any stages", e.ImageName)
}
