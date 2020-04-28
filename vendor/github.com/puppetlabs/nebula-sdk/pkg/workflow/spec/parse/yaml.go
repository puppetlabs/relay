package parse

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

type YAMLTransformer interface {
	Transform(node *yaml.Node) (bool, error)
}

type YAMLDataTransformer struct{}

func (YAMLDataTransformer) Transform(node *yaml.Node) (bool, error) {
	if node.Tag != "!Data" {
		return false, nil
	}

	var query *yaml.Node
	switch node.Kind {
	case yaml.MappingNode:
		if len(node.Content) != 2 || node.Content[0].Value != "query" {
			return false, fmt.Errorf(`expected mapping-style !Data to have exactly one key, "query"`)
		}

		query = node.Content[1]
	case yaml.SequenceNode:
		if len(node.Content) != 1 {
			return false, fmt.Errorf(`expected sequence-style !Data to have exactly one item`)
		}

		query = node.Content[0]
	case yaml.ScalarNode:
		query = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: node.Value,
		}
	}

	// {$type: Data, query: <query>}
	*node = yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "$type"},
			{Kind: yaml.ScalarNode, Value: "Data"},
			{Kind: yaml.ScalarNode, Value: "query"},
			query,
		},
	}
	return true, nil
}

type YAMLSecretTransformer struct{}

func (YAMLSecretTransformer) Transform(node *yaml.Node) (bool, error) {
	if node.Tag != "!Secret" {
		return false, nil
	}

	var name *yaml.Node
	switch node.Kind {
	case yaml.MappingNode:
		if len(node.Content) != 2 || node.Content[0].Value != "name" {
			return false, fmt.Errorf(`expected mapping-style !Secret to have exactly one key, "name"`)
		}

		name = node.Content[1]
	case yaml.SequenceNode:
		if len(node.Content) != 1 {
			return false, fmt.Errorf(`expected sequence-style !Secret to have exactly one item`)
		}

		name = node.Content[0]
	case yaml.ScalarNode:
		name = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: node.Value,
		}
	}

	// {$type: Secret, name: <name>}
	*node = yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "$type"},
			{Kind: yaml.ScalarNode, Value: "Secret"},
			{Kind: yaml.ScalarNode, Value: "name"},
			name,
		},
	}
	return true, nil
}

type YAMLConnectionTransformer struct{}

func (YAMLConnectionTransformer) Transform(node *yaml.Node) (bool, error) {
	if node.ShortTag() != "!Connection" {
		return false, nil
	}

	var connectionType, name *yaml.Node
	switch node.Kind {
	case yaml.MappingNode:
		if len(node.Content) != 4 {
			return false, fmt.Errorf(`expected mapping-style !Connection to have exactly two keys, "type" and "name"`)
		}

		for i := 0; i < len(node.Content); i += 2 {
			switch node.Content[i].Value {
			case "type":
				connectionType = node.Content[i+1]
			case "name":
				name = node.Content[i+1]
			default:
				return false, fmt.Errorf(`expected mapping-style !Connection to have exactly two keys, "type" and "name"`)
			}
		}
	case yaml.SequenceNode:
		if len(node.Content) != 2 {
			return false, fmt.Errorf(`expected mapping-style !Connection to have exactly two items`)
		}

		connectionType = node.Content[0]
		name = node.Content[1]
	default:
		return false, fmt.Errorf(`unexpected scalar value for !Connection, must be a mapping or sequence`)
	}

	// {$type: Connection, type: <type>, name: <name>}
	*node = yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "$type"},
			{Kind: yaml.ScalarNode, Value: "Connection"},
			{Kind: yaml.ScalarNode, Value: "type"},
			connectionType,
			{Kind: yaml.ScalarNode, Value: "name"},
			name,
		},
	}
	return true, nil
}

type YAMLOutputTransformer struct{}

func (YAMLOutputTransformer) Transform(node *yaml.Node) (bool, error) {
	if node.ShortTag() != "!Output" {
		return false, nil
	}

	var from, name *yaml.Node
	switch node.Kind {
	case yaml.MappingNode:
		if len(node.Content) != 4 {
			return false, fmt.Errorf(`expected mapping-style !Output to have exactly two keys, "from" and "name"`)
		}

		for i := 0; i < len(node.Content); i += 2 {
			switch node.Content[i].Value {
			case "from":
				from = node.Content[i+1]
			case "name":
				name = node.Content[i+1]
			default:
				return false, fmt.Errorf(`expected mapping-style !Output to have exactly two keys, "from" and "name"`)
			}
		}
	case yaml.SequenceNode:
		if len(node.Content) != 2 {
			return false, fmt.Errorf(`expected mapping-style !Output to have exactly two items`)
		}

		from = node.Content[0]
		name = node.Content[1]
	default:
		return false, fmt.Errorf(`unexpected scalar value for !Output, must be a mapping or sequence`)
	}

	// {$type: Output, from: <from>, name: <name>}
	*node = yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "$type"},
			{Kind: yaml.ScalarNode, Value: "Output"},
			{Kind: yaml.ScalarNode, Value: "from"},
			from,
			{Kind: yaml.ScalarNode, Value: "name"},
			name,
		},
	}
	return true, nil
}

type YAMLParameterTransformer struct{}

func (YAMLParameterTransformer) Transform(node *yaml.Node) (bool, error) {
	if node.ShortTag() != "!Parameter" {
		return false, nil
	}

	var name *yaml.Node
	switch node.Kind {
	case yaml.MappingNode:
		if len(node.Content) != 2 || node.Content[0].Value != "name" {
			return false, fmt.Errorf(`expected mapping-style !Parameter to have exactly one key, "name"`)
		}

		name = node.Content[1]
	case yaml.SequenceNode:
		if len(node.Content) != 1 {
			return false, fmt.Errorf(`expected sequence-style !Parameter to have exactly one item`)
		}

		name = node.Content[0]
	case yaml.ScalarNode:
		name = &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: node.Value,
		}
	}

	// {$type: Parameter, name: <name>}
	*node = yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "$type"},
			{Kind: yaml.ScalarNode, Value: "Parameter"},
			{Kind: yaml.ScalarNode, Value: "name"},
			name,
		},
	}
	return true, nil
}

type YAMLAnswerTransformer struct{}

func (YAMLAnswerTransformer) Transform(node *yaml.Node) (bool, error) {
	if node.ShortTag() != "!Answer" {
		return false, nil
	}

	var askRef, name *yaml.Node
	switch node.Kind {
	case yaml.MappingNode:
		if len(node.Content) != 4 {
			return false, fmt.Errorf(`expected mapping-style !Answer to have exactly two keys, "askRef" and "name"`)
		}

		for i := 0; i < len(node.Content); i += 2 {
			switch node.Content[i].Value {
			case "askRef":
				askRef = node.Content[i+1]
			case "name":
				name = node.Content[i+1]
			default:
				return false, fmt.Errorf(`expected mapping-style !Answer to have exactly two keys, "askRef" and "name"`)
			}
		}
	case yaml.SequenceNode:
		if len(node.Content) != 2 {
			return false, fmt.Errorf(`expected mapping-style !Answer to have exactly two items`)
		}

		askRef = node.Content[0]
		name = node.Content[1]
	default:
		return false, fmt.Errorf(`unexpected scalar value for !Answer, must be a mapping or sequence`)
	}

	// {$type: Answer, askRef: <askRef>, name: <name>}
	*node = yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "$type"},
			{Kind: yaml.ScalarNode, Value: "Answer"},
			{Kind: yaml.ScalarNode, Value: "askRef"},
			askRef,
			{Kind: yaml.ScalarNode, Value: "name"},
			name,
		},
	}
	return true, nil
}

type YAMLInvocationTransformer struct{}

func (YAMLInvocationTransformer) Transform(node *yaml.Node) (bool, error) {
	tag := node.ShortTag()
	prefix := "!Fn."

	if !strings.HasPrefix(tag, prefix) {
		return false, nil
	}

	name := tag[len(prefix):]
	if len(name) == 0 {
		return false, fmt.Errorf(`expected function name to have the syntax !Fn.<name>`)
	}

	var args *yaml.Node
	switch node.Kind {
	case yaml.MappingNode, yaml.SequenceNode:
		args = &yaml.Node{
			Kind:    node.Kind,
			Content: node.Content,
		}
	case yaml.ScalarNode:
		args = &yaml.Node{
			Kind: yaml.SequenceNode,
			Content: []*yaml.Node{
				{Kind: yaml.ScalarNode, Value: node.Value},
			},
		}
	}

	// {$fn.<name>: <args>}
	*node = yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: fmt.Sprintf("$fn.%s", name)},
			args,
		},
	}
	return true, nil
}

type YAMLBinaryToEncodingTransformer struct{}

func (YAMLBinaryToEncodingTransformer) Transform(node *yaml.Node) (bool, error) {
	if node.ShortTag() != "!!binary" || node.Kind != yaml.ScalarNode {
		return false, nil
	}

	// {$encoding: base64, data: <value>}
	*node = yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "$encoding"},
			{Kind: yaml.ScalarNode, Value: "base64"},
			{Kind: yaml.ScalarNode, Value: "data"},
			{Kind: yaml.ScalarNode, Value: node.Value},
		},
	}
	return true, nil
}

type YAMLUnknownTagTransformer struct{}

func (YAMLUnknownTagTransformer) Transform(node *yaml.Node) (bool, error) {
	if tag := node.ShortTag(); tag != "" && !strings.HasPrefix(tag, "!!") {
		return false, fmt.Errorf(`unknown tag %q`, node.ShortTag())
	}

	return false, nil
}

var YAMLTransformers = []YAMLTransformer{
	YAMLDataTransformer{},
	YAMLSecretTransformer{},
	YAMLConnectionTransformer{},
	YAMLOutputTransformer{},
	YAMLParameterTransformer{},
	YAMLAnswerTransformer{},
	YAMLInvocationTransformer{},
	YAMLBinaryToEncodingTransformer{},
	YAMLUnknownTagTransformer{},
}

func ParseYAML(r io.Reader) (Tree, error) {
	node := &yaml.Node{}
	if err := yaml.NewDecoder(r).Decode(node); err != nil {
		return nil, err
	}

	return ParseYAMLNode(node)
}

func ParseYAMLString(data string) (Tree, error) {
	return ParseYAML(strings.NewReader(data))
}

func ParseYAMLNode(node *yaml.Node) (Tree, error) {
	stack := []*yaml.Node{node}
	for len(stack) > 0 {
		node := stack[0]

		for node.Kind == yaml.AliasNode {
			node = node.Alias
		}

		for _, t := range YAMLTransformers {
			if ok, err := t.Transform(node); err != nil {
				return nil, err
			} else if ok {
				break
			}
		}

		// Remove head and append children for further analysis.
		stack = append(stack[1:], node.Content...)
	}

	var tree interface{}
	if err := node.Decode(&tree); err != nil {
		return nil, err
	}

	return Tree(tree), nil
}
