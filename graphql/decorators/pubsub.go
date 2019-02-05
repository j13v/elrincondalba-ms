package decorators

import (
	"context"

	"github.com/graphql-go/graphql"
)

type contextPubSubKeyType string

const contextPubSubKey contextPubSubKeyType = "pubsub"

/*
ContextRepoConsumer consume the model from graphql context
*/
func ContextPubSubConsumer(handler func(interface{}) (interface{}, error)) func(graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
		return params.Context.Value(contextPubSubKey), nil
	}
}

/*
ContextRepoApply apply the model to graphql context
*/
func ContextPubSubApply(value interface{}) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, contextPubSubKey, value)
	}
}
