package generator

import (
	"os"

	"github.com/puppetlabs/relay/pkg/integration/container/def"
)

const (
	DefaultScriptFilename = "scripts/build-container"
	DefaultRepoNameBase   = "relay.localhost/generated"

	IntermediateRepoNameBase = "relay.localhost/intermediate"
)

type File struct {
	Ref     *def.FileRef
	Mode    os.FileMode
	Content string
}

type Generator struct {
	c *def.Container

	base           *def.FileRef
	scriptFilename string
	repoNameBase   string
}

func (g *Generator) Files() ([]*File, error) {
	var fs []*File

	// Generate build script
	f, err := g.generateScript()
	if err != nil {
		return nil, err
	}
	fs = append(fs, f)

	// Generate Dockerfiles
	ifs, err := g.generateImages()
	if err != nil {
		return nil, err
	}
	fs = append(fs, ifs...)

	return fs, nil
}

type Option func(g *Generator)

func WithScriptFilename(filename string) Option {
	return func(g *Generator) {
		g.scriptFilename = filename
	}
}

func WithRepoNameBase(base string) Option {
	return func(g *Generator) {
		g.repoNameBase = base
	}
}

func WithFilesRelativeTo(ref *def.FileRef) Option {
	return func(g *Generator) {
		g.base = ref.Dir()
	}
}

func New(container *def.Container, opts ...Option) *Generator {
	g := &Generator{
		c: container,

		base:           def.NewFileRef("."),
		scriptFilename: DefaultScriptFilename,
		repoNameBase:   DefaultRepoNameBase,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}
