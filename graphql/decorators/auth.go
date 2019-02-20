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
// type User struct {
// 	Name   string
// 	Scopes []string
// }

// type contextUserKeyType string

// const contextUserKey contextUserKeyType = "user"

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

	// secret := []byte("secret")
	// 		email := "coco@cookies.com"
	// 		password := "I love cupcakes"
	// 		argEmail := params.Args["email"].(string)
	// 		argSignarute := params.Args["signature"].(string)

	// 		mac := hmac.New(sha256.New, []byte(secret))
	// 		mac.Write([]byte(fmt.Sprintf("%s:%s", email, password)))
	// 		phash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	// 		// Create a new token object, specifying signing method and the claims
	// 		// you would like it to contain.
	// 		if argEmail == email && argSignarute == string(phash) {
	// 			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 				"admin": true,
	// 				"nbf":   time.Now().Unix(),
	// 				"exp":   time.Now().Add(time.Hour * 72).Unix(),
	// 			})

	// 			// Sign and get the complete encoded token as a string using the secret
	// 			tokenString, err := token.SignedString(secret)
	// 			return tokenString, err
	// 		}

	return func(ctx context.Context) context.Context {
		fmt.Printf("%v", ctx)
		// return context.WithValue(ctx, contextUserKey, user)
		return ctx
	}
}

// for _, accountMock := range accountsMock {
// 	if accountMock.Username == account.(User).Username {
// 		return accountMock, nil
// 	}
// }
// return &User{}, nil

// if user, userOk := params.Context.Value(contextUserKey).(User); userOk {
// 	userRolesOk, missingRole := util.Contains(scopes, user.Scopes, true)
// 	if userRolesOk {
// 		res, err := handler(params)
// 		return res, err
// 	}
// 	return nil, fmt.Errorf("Forbidden user %s has not required scope %s", user.Name, missingRole)

// }
// return nil, errors.New("Missing user context")
