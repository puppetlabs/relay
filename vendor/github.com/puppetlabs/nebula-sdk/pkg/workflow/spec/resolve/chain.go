package resolve

import (
	"context"
)

type chainDataTypeResolvers struct {
	resolvers []DataTypeResolver
}

func (cr *chainDataTypeResolvers) ResolveData(ctx context.Context, query string) (interface{}, error) {
	for _, r := range cr.resolvers {
		s, err := r.ResolveData(ctx, query)
		if _, ok := err.(*DataNotFoundError); ok {
			continue
		} else if err != nil {
			return "", err
		}

		return s, nil
	}

	return "", &DataNotFoundError{Query: query}
}

func ChainDataTypeResolvers(resolvers ...DataTypeResolver) DataTypeResolver {
	return &chainDataTypeResolvers{resolvers: resolvers}
}

type chainSecretTypeResolvers struct {
	resolvers []SecretTypeResolver
}

func (cr *chainSecretTypeResolvers) ResolveSecret(ctx context.Context, name string) (string, error) {
	for _, r := range cr.resolvers {
		s, err := r.ResolveSecret(ctx, name)
		if _, ok := err.(*SecretNotFoundError); ok {
			continue
		} else if err != nil {
			return "", err
		}

		return s, nil
	}

	return "", &SecretNotFoundError{Name: name}
}

func ChainSecretTypeResolvers(resolvers ...SecretTypeResolver) SecretTypeResolver {
	return &chainSecretTypeResolvers{resolvers: resolvers}
}

type chainConnectionTypeResolvers struct {
	resolvers []ConnectionTypeResolver
}

func (cr *chainConnectionTypeResolvers) ResolveConnection(ctx context.Context, connectionType, name string) (interface{}, error) {
	for _, r := range cr.resolvers {
		o, err := r.ResolveConnection(ctx, connectionType, name)
		if _, ok := err.(*ConnectionNotFoundError); ok {
			continue
		} else if err != nil {
			return "", err
		}

		return o, nil
	}

	return "", &ConnectionNotFoundError{Type: connectionType, Name: name}
}

func ChainConnectionTypeResolvers(resolvers ...ConnectionTypeResolver) ConnectionTypeResolver {
	return &chainConnectionTypeResolvers{resolvers: resolvers}
}

type chainOutputTypeResolvers struct {
	resolvers []OutputTypeResolver
}

func (cr *chainOutputTypeResolvers) ResolveOutput(ctx context.Context, from, name string) (interface{}, error) {
	for _, r := range cr.resolvers {
		o, err := r.ResolveOutput(ctx, from, name)
		if _, ok := err.(*OutputNotFoundError); ok {
			continue
		} else if err != nil {
			return "", err
		}

		return o, nil
	}

	return "", &OutputNotFoundError{From: from, Name: name}
}

func ChainOutputTypeResolvers(resolvers ...OutputTypeResolver) OutputTypeResolver {
	return &chainOutputTypeResolvers{resolvers: resolvers}
}

type chainParameterTypeResolvers struct {
	resolvers []ParameterTypeResolver
}

func (cr *chainParameterTypeResolvers) ResolveParameter(ctx context.Context, name string) (interface{}, error) {
	for _, r := range cr.resolvers {
		p, err := r.ResolveParameter(ctx, name)
		if _, ok := err.(*ParameterNotFoundError); ok {
			continue
		} else if err != nil {
			return nil, err
		}

		return p, nil
	}

	return nil, &ParameterNotFoundError{Name: name}
}

func ChainParameterTypeResolvers(resolvers ...ParameterTypeResolver) ParameterTypeResolver {
	return &chainParameterTypeResolvers{resolvers: resolvers}
}

type chainAnswerTypeResolvers struct {
	resolvers []AnswerTypeResolver
}

func (cr *chainAnswerTypeResolvers) ResolveAnswer(ctx context.Context, askRef, name string) (interface{}, error) {
	for _, r := range cr.resolvers {
		p, err := r.ResolveAnswer(ctx, askRef, name)
		if _, ok := err.(*AnswerNotFoundError); ok {
			continue
		} else if err != nil {
			return nil, err
		}

		return p, nil
	}

	return nil, &AnswerNotFoundError{AskRef: askRef, Name: name}
}

func ChainAnswerTypeResolvers(resolvers ...AnswerTypeResolver) AnswerTypeResolver {
	return &chainAnswerTypeResolvers{resolvers: resolvers}
}
