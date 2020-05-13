package def

import (
	"io/ioutil"
	"net/http"

	"github.com/imdario/mergo"
	v1 "github.com/puppetlabs/relay/pkg/integration/container/types/v1"
)

type Image struct {
	// TemplateName is the file name of the template.
	TemplateName string

	// TemplateData is the file content of the template.
	TemplateData string

	// DependsOn is the graph showing the hierarchy of other images required to
	// build the template.
	DependsOn []string
}

type Setting struct {
	Description string
	Value       interface{}
}

type Common struct {
	SDKVersion string
	Images     map[string]*Image
	Settings   map[string]*Setting

	resolver *Resolver
}

type CommonOption func(c *Common)

func WithResolver(resolver *Resolver) CommonOption {
	return func(c *Common) {
		c.resolver = resolver
	}
}

func WithSetting(name string, value interface{}) CommonOption {
	return func(c *Common) {
		if c.Settings[name] == nil {
			c.Settings[name] = &Setting{}
		}

		c.Settings[name].Value = value
	}
}

func NewCommonFromTyped(sctt *v1.StepContainerCommon, opts ...CommonOption) (*Common, error) {
	c := &Common{
		SDKVersion: sctt.SDKVersion,
		Settings:   make(map[string]*Setting, len(sctt.Settings)),
		Images:     make(map[string]*Image, len(sctt.Images)),
	}

	// Check that the SDK version has been set to a sane value.
	switch c.SDKVersion {
	case "v1": // TODO: Should this be formalized somewhere?
	case "":
		c.SDKVersion = "v1"
	default:
		return nil, &UnknownSDKVersionError{Got: c.SDKVersion}
	}

	// Set a base resolver so that we always have one defined.
	WithResolver(DefaultResolver)(c)

	// Add settings from the typed definition.
	for name, setting := range sctt.Settings {
		c.Settings[name] = &Setting{
			Description: setting.Description,
			Value:       setting.Value,
		}
	}

	// Apply options before processing further.
	for _, opt := range opts {
		opt(c)
	}

	for name, image := range sctt.Images {
		ci := &Image{
			TemplateName: image.Template.Name,
			DependsOn:    image.DependsOn,
		}

		fr, err := NewFileRefFromTyped(image.Template, WithFileRefResolver(c.resolver))
		if err != nil {
			return nil, err
		}

		if err := fr.WithFile(func(f http.File) error {
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}

			ci.TemplateData = string(b)
			return nil
		}); err != nil {
			return nil, err
		}

		c.Images[name] = ci
	}

	// We need to merge this file with its parents.
	if sctt.Inherit != nil {
		fr, err := NewFileRefFromTyped(*sctt.Inherit, WithFileRefResolver(c.resolver))
		if err != nil {
			return nil, err
		}

		parent, err := NewTemplateFromFileRef(fr)
		if err != nil {
			return nil, &TemplateError{FileRef: fr, Cause: err}
		}

		if err := mergo.Merge(c, parent.Common); err != nil {
			return nil, &TemplateError{FileRef: fr, Cause: err}
		}
	}

	return c, nil
}
