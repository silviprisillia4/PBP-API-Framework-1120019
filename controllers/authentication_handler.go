package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var jwtKey = []byte("rahasia123")
var tokenName = "token"

type Claims struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserType int    `json:"type"`
	jwt.StandardClaims
}

func generateToken(c echo.Context, id int, name string, userType int) {
	tokenExpiryTime := time.Now().Add(5 * time.Minute) // expirednya 5 menit

	// create claims with user data
	claims := &Claims{
		ID:       id,
		Name:     name,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}

	// encrpyt claim to jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // claims to jwt
	signedToken, err := token.SignedString(jwtKey)             // jwt key
	if err != nil {
		return
	}

	c.SetCookie(&http.Cookie{
		Name:     tokenName,
		Value:    signedToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
	})

}

func resetUserToken(c echo.Context) {
	// reset cookie (logout)
	c.SetCookie(&http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
}

func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessType := 0
		isValidToken := validateUserToken(c, accessType)
		if !isValidToken {
			return sendUnauthorizedResponse(c)
		} else {
			return next(c)
		}
	}
}

func validateUserToken(c echo.Context, accessType int) bool {
	isAccessTokenValid, id, name, userType := validateTokenFromCookies(c)
	fmt.Println(id, "\t", name, "\t", userType, "\t", accessType, "\t", isAccessTokenValid)
	if isAccessTokenValid {
		isUserValid := userType == accessType // ngecek typenya kita sama type dari token sama atau ga
		if isUserValid {
			return true
		}
	}
	return false
}

func validateTokenFromCookies(c echo.Context) (bool, int, string, int) {
	if cookie, err := c.Cookie(tokenName); err == nil {
		accessToken := cookie.Value // karena token disimpan dalam atribut Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(accessToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.ID, accessClaims.Name, accessClaims.UserType
		}
	}
	return false, -1, "", -1
}
