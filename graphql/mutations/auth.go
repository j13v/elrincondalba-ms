package mutations

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
)

var MutationAuth = graphql.Fields{
	"getAccessToken": &graphql.Field{
		Type: graphql.String,
		Args: graphql.FieldConfigArgument{
			"signature": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Description: "List catalog",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			secret := []byte(params.Context.Value("secret").(string))
			email := os.Getenv("API_ADMIN_EMAIL")
			password := os.Getenv("API_ADMIN_PASSWORD")
			argSignarute := params.Args["signature"].(string)

			mac := hmac.New(sha256.New, secret)
			mac.Write([]byte(fmt.Sprintf("%s:%s", email, password)))
			phash := base64.StdEncoding.EncodeToString(mac.Sum(nil))
			fmt.Printf("%v %v\n", string(phash), argSignarute)
			// Create a new token object, specifying signing method and the claims
			// you would like it to contain.
			if argSignarute == string(phash) {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"admin": true,
					"nbf":   time.Now().Unix(),
					"exp":   time.Now().Add(time.Hour * 72).Unix(),
				})
				// Sign and get the complete encoded token as a string using the secret
				tokenString, err := token.SignedString(secret)
				fmt.Print(tokenString);
				return tokenString, err
			}

			return nil, fmt.Errorf("Invalid user login for %s", email)
		},
	},
}
