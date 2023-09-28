package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CreateToken Create  jwt token with 'authorized' and 'user_id' claims
func CreateToken(
	user_id uint32,
	tokenExpireSeconds uint,
	apiSecret string,
) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Second * time.Duration(tokenExpireSeconds)).Unix() //Token expires after 'tokenExpireSeconds' seconds
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(apiSecret))
}

func ExtractTokenID(
	r *http.Request,
	apiSecret string,
) (uint32, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(apiSecret), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint32(uid), nil
	}
	return 0, nil
}

// extractToken gets token from any url header or cookie
func extractToken(r *http.Request) string {
	// from url
	keys := r.URL.Query()
	token := keys.Get("access_token")
	if token != "" {
		return token
	}

	// from header
	const sepTypeR = " "
	var bearerToken = r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, sepTypeR)) == 2 {
		return strings.Split(bearerToken, sepTypeR)[1]
	}

	if bearerToken != "" {
		return bearerToken
	}
	// from cookie
	var tokenCookie, _ = r.Cookie("auth")
	if tokenCookie != nil {
		return tokenCookie.Value
	}

	return ""

}
