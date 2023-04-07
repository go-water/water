package water

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"net/http"
	"strings"
	"time"
)

func SetAuthToken(user, key string, expire time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    user,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
	}

	signingKey := []byte(key)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", jwt.ErrTokenSignatureInvalid
	}

	return tokenString, nil
}

func Valid(req *http.Request, key string) (*jwt.RegisteredClaims, error) {
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}, request.WithClaims(&jwt.RegisteredClaims{}))
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenSignatureInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
