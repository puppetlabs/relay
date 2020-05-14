package sdk

import (
	"net/http"
	"sync"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/puppetlabs/relay/pkg/debug"
)

// Clone uses the git CLI to pull down the SDK repo for usage locally.
func Clone() (billy.Filesystem, error) {
	fs := memfs.New()
	store := memory.NewStorage()

	_, err := git.Clone(store, fs, &git.CloneOptions{
		URL:      "https://github.com/relay-integrations/container-definitions",
		Progress: debug.Writer(),
	})

	if err != nil {
		return nil, err
	}

	return fs, nil
}

type gitFileSystem struct {
	fs     billy.Filesystem
	cloner sync.Once
}

func (fs *gitFileSystem) doClone() (err error) {
	fs.cloner.Do(func() {
		fs.fs, err = Clone()
	})

	return
}

func (fs *gitFileSystem) Open(path string) (http.File, error) {
	if err := fs.doClone(); err != nil {
		return nil, err
	}

	// we have to stat the file first because we need to see if it's a directory.
	// if it is then we have to treat it slightly differently because Billy
	// doesn't support "opening" directories.
	//
	// we likewise re-use the stat data further down.
	info, err := fs.fs.Stat(path)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		f, err := fs.fs.Open(path)

		if err != nil {
			return nil, err
		}

		return billyFile{f, info}, nil
	} else {
		return billyFile{stat: info}, nil
	}
}

func NewFileSystem() http.FileSystem {
	return &gitFileSystem{}
}
