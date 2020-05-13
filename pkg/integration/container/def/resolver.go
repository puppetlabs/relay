package def

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/integration/container/asset"
)

// Resolver allows the generator to load dependent resources.
type Resolver struct {
	// FileSystem is the filesystem implementation to use to load dependent
	// resources.
	FileSystem http.FileSystem

	// WorkingDirectory is the directory to use to resolve relative paths in the
	// step container data.
	WorkingDirectory string
}

func (r *Resolver) Local() bool {
	return r == nil || r.FileSystem == nil
}

func (r *Resolver) Open(name string) (http.File, error) {
	if r == nil {
		r = DefaultResolver
	}

	if !path.IsAbs(name) {
		name = path.Join(r.WorkingDirectory, name)
	}

	if !r.Local() {
		return r.FileSystem.Open(name)
	}

	return os.Open(filepath.FromSlash(name))
}

func (r *Resolver) NewRelativeTo(name string) *Resolver {
	if r == nil {
		r = DefaultResolver
	}

	dir := path.Dir(name)

	if !path.IsAbs(name) {
		dir = path.Join(r.WorkingDirectory, dir)
	}

	return &Resolver{
		FileSystem:       r.FileSystem,
		WorkingDirectory: dir,
	}
}

func (r *Resolver) String() string {
	if r == nil {
		r = DefaultResolver
	}

	var fs string

	switch r {
	case DefaultResolver:
	case SDKResolver:
		fs = "sdk:"
	default:
		fs = "(unknown):"
	}

	return fs + r.WorkingDirectory
}

var (
	DefaultResolver = &Resolver{}
	SDKResolver     = &Resolver{
		FileSystem:       asset.FileSystem,
		WorkingDirectory: "templates",
	}
)
