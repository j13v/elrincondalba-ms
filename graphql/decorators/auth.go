package decorators

import (
	"context"
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/mitchellh/mapstructure"
)

// https://gist.github.com/cryptix/45c33ecf0ae54828e63b
// https://www.thepolyglotdeveloper.com/2018/07/jwt-authorization-graphql-api-using-golang/

func ValidateJWT(t string, secret string) (interface{}, error) {

	jwtSecret := []byte(secret)

	if t == "" {
		return nil, errors.New("Authorization token must be present")
	}
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return jwtSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var decodedToken interface{}
		mapstructure.Decode(claims, &decodedToken)
		return decodedToken, nil
	} else {
		return nil, errors.New("Invalid authorization token")
	}
}

/*
ContextUserScopesValidator consume the user from graphql context
*/
/// scopes []string,
func ContextAuthConsumer(handler func(graphql.ResolveParams, interface{}) (interface{}, error)) func(graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
		claims, err := ValidateJWT(
			params.Context.Value("token").(string),
			params.Context.Value("secret").(string),
		)

		if err == nil {
			res, err := handler(params, claims)
			return res, err
		}
		return nil, err
	}
}

func ContextAuthIsAdmin(handler func(graphql.ResolveParams) (interface{}, error)) func(graphql.ResolveParams) (interface{}, error) {
	return ContextAuthConsumer(func(params graphql.ResolveParams, claims interface{}) (interface{}, error) {
		mapClaims := claims.(jwt.MapClaims)

		if mapClaims["admin"] == true {
			res, err := handler(params)
			return res, err
		}

		return nil, errors.New("Forbidden")
	})
}

/*
ContextUserApply consume the user from graphql context
*/
func ContextAuthApply() func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return ctx
	}
}