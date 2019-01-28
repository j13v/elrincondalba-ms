package decorators

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/mongodb"
)

type contextModelKeyType string

const contextModelKey contextModelKeyType = "model"

/*
ContextModelConsumer consume the model from graphql context
*/
func ContextModelConsumer(handler func(graphql.ResolveParams, mongodb.Model) (interface{}, error)) func(graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
		if model, ok := params.Context.Value("model").(mongodb.Model); ok {
			res, err := handler(params, model)
			// fmt.Printf("%+v\n", ok)
			return res, err
		}
		return nil, errors.New("Model context not found")
	}
}

/*
ContextModelApply apply the model to graphql context
*/
func ContextModelApply(model mongodb.Model) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "model", model)
	}
}
