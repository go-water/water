package water

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

func SetAuthToken(uniqueUser, issuer, privateKeyPath string, expire time.Duration) (tokenString string, err error) {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	claims := jwt.RegisteredClaims{
		ID:        uniqueUser,
		Issuer:    issuer,
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

// ParseFromRequest 兼容 http,ws
func ParseFromRequest(req *http.Request, publicKeyPath string) (uniqueUser, issuer, signature string, err error) {
	token, err := request.ParseFromRequest(req, request.AuthorizationHeaderExtractor, func(t *jwt.Token) (any, error) {
		publicKey, innErr := os.ReadFile(publicKeyPath)
		if innErr != nil {
			return "", innErr
		}

		return jwt.ParseRSAPublicKeyFromPEM(publicKey)
	}, request.WithClaims(&jwt.RegisteredClaims{}))

	if token != nil && token.Valid {
		return parseToken(token)
	}

	// 兼容 ws
	if wsp := req.Header.Get("Sec-Websocket-Protocol"); len(wsp) > 0 {
		token, err = jwt.ParseWithClaims(wsp, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			publicKey, innErr := os.ReadFile(publicKeyPath)
			if innErr != nil {
				return "", innErr
			}

			return jwt.ParseRSAPublicKeyFromPEM(publicKey)
		})
		if err != nil {
			return "", "", "", err
		}

		if token.Valid {
			return parseToken(token)
		}
	}

	return "", "", "", jwt.ErrTokenSignatureInvalid
}

func parseToken(token *jwt.Token) (uniqueUser, issuer, signature string, err error) {
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", "", "", jwt.ErrTokenInvalidClaims
	}

	return claims.ID, claims.Issuer, string(token.Signature), nil
}
