package sdk

import (
	"net/http"

	"github.com/puppetlabs/relay/pkg/debug"
)

type githubFileSystem struct {
	c githubClient
}

func (fs githubFileSystem) Open(name string) (http.File, error) {
	content, err := fs.c.Get(name)

	if err != nil {
		debug.Logf("error fetching content from SDK: %v", err.Error())
		return nil, err
	}

	return newGithubFile(name, content, &fs.c)
}

func NewFileSystem() http.FileSystem {
	return githubFileSystem{}
}
