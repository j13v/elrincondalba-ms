package decorators

import (
	"context"
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/jal88/elrincondalba-ms/util"
)

type User struct {
	Name   string
	Scopes []string
}
type contextUserKeyType string

const contextUserKey contextUserKeyType = "user"

/*
ContextUserScopesValidator consume the user from graphql context
*/
func ContextUserScopesValidator(scopes []string, handler func(graphql.ResolveParams) (interface{}, error)) func(graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
		if user, userOk := params.Context.Value(contextUserKey).(User); userOk {
			userRolesOk, missingRole := util.Contains(scopes, user.Scopes, true)
			if userRolesOk {
				res, err := handler(params)
				return res, err
			}
			return nil, fmt.Errorf("Forbidden user %s has not required scope %s", user.Name, missingRole)

		}
		return nil, errors.New("Missing user context")
	}
}

/*
ContextUserApply consume the user from graphql context
*/
func ContextUserApply(user User) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, contextUserKey, user)
	}
}
