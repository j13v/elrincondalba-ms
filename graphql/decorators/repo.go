package decorators

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

type contextRepoKeyType string

const contextRepoKey contextRepoKeyType = "repo"

/*
ContextRepoConsumer consume the model from graphql context
*/
func ContextRepoConsumer(handler func(graphql.ResolveParams, mongodb.Repo) (interface{}, error)) func(graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
		if model, ok := params.Context.Value(contextRepoKey).(mongodb.Repo); ok {
			res, err := handler(params, model)
			return res, err
		}
		return nil, errors.New("Repo context not found")
	}
}

/*
ContextRepoApply apply the model to graphql context
*/
func ContextRepoApply(model mongodb.Repo) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, contextRepoKey, model)
	}
}
