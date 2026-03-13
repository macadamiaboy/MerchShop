package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/macadamiaboy/AvitoMerchShop/internal/config"
)

var secretSalt = []byte(config.LoadConfigData().Auth.SecretSalt)

func GenToken(id int64) (string, error) {
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		Issuer:    "MerchShop",
		Subject:   strconv.FormatInt(id, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(secretSalt)
	return s, err
}

func getClaims(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretSalt, nil
	})

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

/*
func Verify(tokenString string, login string) error {
	claims, err := getClaims(tokenString)
	if err != nil {
		return err
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("The token has expired")
	}

	if claims.Subject != login {
		return fmt.Errorf("The token is incorrect")
	}

	return nil
}
*/

func GetIdFromToken(tokenString string) (int64, error) {
	claims, err := getClaims(tokenString)
	if err != nil {
		return 0, err
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return 0, fmt.Errorf("The token has expired")
	}

	res, err := strconv.ParseInt(claims.Subject, 10, 64)

	return res, nil
}
