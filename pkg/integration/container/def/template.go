package def

import (
	"io"
	"net/http"

	v1 "github.com/puppetlabs/relay/pkg/integration/container/types/v1"
)

type Template struct {
	*Common
}

func NewTemplateFromTyped(sctt *v1.StepContainerTemplate, opts ...CommonOption) (*Template, error) {
	co, err := NewCommonFromTyped(sctt.StepContainerCommon, opts...)
	if err != nil {
		return nil, err
	}

	t := &Template{
		Common: co,
	}
	return t, nil
}

func NewTemplateFromReader(r io.Reader, opts ...CommonOption) (*Template, error) {
	sctt, err := v1.NewStepContainerTemplateFromReader(r)
	if err != nil {
		return nil, err
	}

	return NewTemplateFromTyped(sctt, opts...)
}

func NewTemplateFromFileRef(ref *FileRef) (t *Template, err error) {
	err = ref.WithFile(func(f http.File) (err error) {
		fi, err := f.Stat()
		if err != nil {
			return err
		} else if fi.IsDir() {
			t, err = NewTemplateFromFileRef(ref.Join(DefaultFilename))
		} else {
			t, err = NewTemplateFromReader(f, WithResolver(ref.ResolverHere()))
		}
		return
	})
	return
}
