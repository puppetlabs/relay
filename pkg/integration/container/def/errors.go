package def

import (
	"errors"
	"fmt"

	v1 "github.com/puppetlabs/relay/pkg/integration/container/types/v1"
)

var (
	ErrMissingName = errors.New("def: could not automatically determine container name from parent directory")
)

type UnknownFileSourceError struct {
	Got v1.FileSource
}

func (e *UnknownFileSourceError) Error() string {
	return fmt.Sprintf("def: unknown file source: %q", e.Got)
}

type UnknownSDKVersionError struct {
	Got string
}

func (e *UnknownSDKVersionError) Error() string {
	return fmt.Sprintf("def: unknown SDK version %q", e.Got)
}

type MissingSettingValueError struct {
	Name string
}

func (e *MissingSettingValueError) Error() string {
	return fmt.Sprintf("def: setting %q has no value", e.Name)
}

type TemplateError struct {
	FileRef *FileRef
	Cause   error
}

func (e *TemplateError) Error() string {
	return fmt.Sprintf("def: error in parent template %q: %+v", e.FileRef, e.Cause)
}
