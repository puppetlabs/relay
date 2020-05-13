package def

import (
	"net/http"
	"path"
	"path/filepath"

	v1 "github.com/puppetlabs/relay/pkg/integration/container/types/v1"
)

// FileRef points to a given file in a resolver.
type FileRef struct {
	resolver *Resolver
	name     string
}

func (fr *FileRef) Local() bool {
	return fr.resolver.Local()
}

func (fr *FileRef) Dir() *FileRef {
	return NewFileRef(path.Dir(fr.name), WithFileRefResolver(fr.resolver))
}

func (fr *FileRef) Join(name string) *FileRef {
	return NewFileRef(path.Join(fr.name, name), WithFileRefResolver(fr.resolver))
}

// ResolverHere returns a resolver that can look up files relative to this
// FileRef.
func (fr *FileRef) ResolverHere() *Resolver {
	return fr.resolver.NewRelativeTo(fr.name)
}

func (fr *FileRef) WithFile(fn func(f http.File) error) error {
	f, err := fr.resolver.Open(fr.name)
	if err != nil {
		return err
	}
	defer f.Close()

	return fn(f)
}

func (fr *FileRef) Name() string {
	return fr.name
}

func (fr *FileRef) String() string {
	if fr.resolver == nil {
		return fr.name
	}

	return path.Join(fr.resolver.String(), fr.name)
}

type FileRefOption func(fr *FileRef)

func WithFileRefResolver(resolver *Resolver) FileRefOption {
	return func(fr *FileRef) {
		fr.resolver = resolver
	}
}

func NewFileRef(name string, opts ...FileRefOption) *FileRef {
	fr := &FileRef{name: name}
	for _, opt := range opts {
		opt(fr)
	}

	if fr.resolver == nil || fr.resolver.FileSystem == nil {
		fr.name = filepath.ToSlash(fr.name)
	}

	return fr
}

func NewFileRefFromTyped(ref v1.FileRef, opts ...FileRefOption) (*FileRef, error) {
	switch ref.From {
	case v1.FileSourceSystem:
		return NewFileRef(ref.Name, opts...), nil
	case v1.FileSourceSDK:
		// Clean name so users can't traverse outside of the template directory
		// and override the resolver with the SDK one.
		name := path.Clean(`.` + path.Clean(`/`+ref.Name))
		opts = append([]FileRefOption{}, opts...)
		opts = append(opts, WithFileRefResolver(SDKResolver))
		return NewFileRef(name, opts...), nil
	default:
		return nil, &UnknownFileSourceError{Got: ref.From}
	}
}
