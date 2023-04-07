package water

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"net/http"
	"strings"
	"time"
)

func SetAuthToken(uniqueUser, privateKey string, expire time.Duration) (tokenString string, err error) {
	claims := jwt.RegisteredClaims{
		Issuer:    uniqueUser,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
	}

	signingKey := []byte(privateKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err = token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", jwt.ErrTokenSignatureInvalid
	}

	return tokenString, nil
}

func ParseAndValid(req *http.Request, privateKey string) (uniqueUser string, err error) {
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(t *jwt.Token) (interface{}, error) {
		return []byte(privateKey), nil
	}, request.WithClaims(&jwt.RegisteredClaims{}))
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrTokenSignatureInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", jwt.ErrTokenInvalidClaims
	}

	return claims.Issuer, nil
}
