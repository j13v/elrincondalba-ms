package queries

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
)

var QueryAuth = graphql.Fields{
	"getAuthToken": &graphql.Field{
		Type: graphql.String,
		Args: graphql.FieldConfigArgument{
			"email": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"signature": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Description: "List catalog",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			secret := []byte("secret")
			email := "coco@cookies.com"
			password := "I love cupcakes"
			argEmail := params.Args["email"].(string)
			argSignarute := params.Args["signature"].(string)

			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write([]byte(fmt.Sprintf("%s:%s", email, password)))
			phash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
			// Create a new token object, specifying signing method and the claims
			// you would like it to contain.
			if argEmail == email && argSignarute == string(phash) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"admin": true,
					"nbf":   time.Now().Unix(),
					"exp":   time.Now().Add(time.Hour * 72).Unix(),
				})

				// Sign and get the complete encoded token as a string using the secret
				tokenString, err := token.SignedString(secret)
				return tokenString, err
			}

			return nil, fmt.Errorf("Invalid user login for %s", email)
		},
	},
}
