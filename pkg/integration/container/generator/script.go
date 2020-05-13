package generator

import (
	"bytes"
	"path"
	"sort"
	"text/template"

	"github.com/puppetlabs/horsehead/v2/graph"
	"github.com/puppetlabs/horsehead/v2/graph/algo"
	"github.com/puppetlabs/horsehead/v2/graph/traverse"
	"github.com/puppetlabs/relay/pkg/integration/container/asset"
)

var (
	scriptTemplate = template.Must(template.New("build.tpl").Parse(asset.MustAssetString("scripts/build.tpl")))
)

type scriptTemplateImageData struct {
	// Name is the final repository name of the image.
	Name string

	// Ref is the internal intermediate build reference name of the image.
	Ref string

	// TagSuffix is the additional suffix, if any, associated with this image.
	TagSuffix string

	// Filename is the path to the Dockerfile for the image.
	Filename string
}

type scriptTemplateData struct {
	Images []*scriptTemplateImageData
}

func (gen *Generator) generateScript() (*File, error) {
	ref := gen.base.Join(gen.scriptFilename)

	// Organize the dependencies into a graph to determine which should come
	// first.
	g := graph.NewSimpleDirectedGraphWithFeatures(graph.DeterministicIteration)

	// Get the names of the images and sort them for consistency.
	imageNames := make([]string, len(gen.c.Images))
	i := 0
	for name := range gen.c.Images {
		imageNames[i] = name
		i++
	}

	sort.Strings(imageNames)

	// For each image, add it as a vertex.
	for _, name := range imageNames {
		g.AddVertex(name)
	}

	// Now connect its dependencies.
	for name, image := range gen.c.Images {
		for _, dep := range image.DependsOn {
			if err := g.Connect(dep, name); err == graph.ErrEdgeAlreadyInGraph {
				continue
			} else if err == graph.ErrWouldCreateLoop {
				return nil, &ImageDependencySelfReferenceError{ImageName: name}
			} else if _, ok := err.(*graph.VertexNotFoundError); ok {
				return nil, &ImageDependencyMissingError{ImageName: name, Want: dep}
			} else if err != nil {
				return nil, err
			}
		}
	}

	// Cycle check.
	var cycles [][]string
	algo.TiernanSimpleCyclesOf(g).CyclesInto(&cycles)
	if len(cycles) > 0 {
		return nil, &ImageDependencyCyclesError{Cycles: cycles}
	}

	// Write out the build instructions in topological order.
	data := &scriptTemplateData{
		Images: make([]*scriptTemplateImageData, len(gen.c.Images)),
	}

	i = 0
	traverse.NewTopologicalOrderTraverser(g).ForEachInto(func(imageName string) error {
		data.Images[i] = &scriptTemplateImageData{
			Name:      path.Join(gen.repoNameBase, gen.c.Name),
			Ref:       gen.imageRef(imageName),
			TagSuffix: imageTagSuffix(imageName),
			Filename:  dockerfile(imageName).Filename(),
		}
		i++

		return nil
	})

	var buf bytes.Buffer
	if err := scriptTemplate.Execute(&buf, data); err != nil {
		// This should really never happen.
		return nil, err
	}

	f := &File{
		Ref:     ref,
		Mode:    0755,
		Content: buf.String(),
	}
	return f, nil
}
