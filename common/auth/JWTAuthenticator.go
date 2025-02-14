package auth

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/libileh/eegis/common/errors"
	"github.com/libileh/eegis/common/properties"
)

type JWTAuth struct {
	secretKey string
	audience  string
	issuer    string
}

type JWTAuthenticator struct {
	JwtAuth    *JWTAuth
	HttpError  *errors.Error
	Properties *properties.AuthProperties
	Validator  *validator.Validate
}

func NewJWTAuth(secretKey, audience, issuer string) *JWTAuth {
	return &JWTAuth{secretKey: secretKey, audience: audience, issuer: issuer}
}

func NewJWTAuthenticator(Properties *properties.AuthProperties, HttpError *errors.Error,
	Validator *validator.Validate) *JWTAuthenticator {
	return &JWTAuthenticator{
		JwtAuth:    NewJWTAuth(Properties.Token.Secret, Properties.Audience, Properties.Issuer),
		HttpError:  HttpError,
		Properties: Properties,
		Validator:  Validator,
	}
}

func (auth *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(auth.JwtAuth.secretKey))
}

func (auth *JWTAuthenticator) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(auth.JwtAuth.secretKey), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(auth.JwtAuth.audience),
		jwt.WithIssuer(auth.JwtAuth.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
