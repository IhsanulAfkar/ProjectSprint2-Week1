package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"
	"week1/models"

	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	jwt.RegisteredClaims
	PkId 	int `json:"pkId"`
	Email    string `json:"email"`
}

func SignJWT(user models.User) string {
	// expiredIn := 28800 // 8 hours
	exp := time.Now().Add(time.Hour * 8)
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer: "Cat Socials",
		},
		PkId: user.PkId,
		Email: user.Email,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return signedToken

}
func ParseToken(jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	parsedToken, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		return "", errors.New("unable to parse claims")
	}
	id:=fmt.Sprint(parsedToken["pkId"])
	fmt.Println(id)
	return id, nil
}