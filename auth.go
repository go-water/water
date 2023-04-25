package water

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"net/http"
	"strings"
	"time"
)

func SetAuthToken(uniqueUser, privateKey string, expire time.Duration) (tokenString string, err error) {
	pri, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(pri)
	if err != nil {
		return "", err
	}

	claims := jwt.RegisteredClaims{
		Issuer:    uniqueUser,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
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

func ParseAndValid(req *http.Request, publicKey string) (uniqueUser, signature string, err error) {
	pub, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", "", err
	}

	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(t *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(pub)
	}, request.WithClaims(&jwt.RegisteredClaims{}))
	if err != nil {
		return "", "", err
	}

	if !token.Valid {
		return "", "", jwt.ErrTokenSignatureInvalid
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", "", jwt.ErrTokenInvalidClaims
	}

	return claims.Issuer, token.Signature, nil
}
