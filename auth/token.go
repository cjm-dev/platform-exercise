package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"strings"
	"time"
)

type UserAuth struct {
	UserID 		int
	UserUUID	string
}
// todo - implement refresh tokens
func CreateToken(ua UserAuth) (string, error) {
	expireTime, err := time.ParseDuration(os.Getenv("TIMEOUT"))
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["auth_uuid"] = ua.UserUUID
	claims["user_id"] = ua.UserID
	claims["exp"] = time.Now().Add(time.Minute * expireTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func RetrieveJwtToken(bt string) string {
	var token string
	tokenData := strings.Split(bt, " ")
	if len(tokenData) == 2 {
		token = tokenData[1]
	}
	if token != "" {
		return token
	}
	return ""
}

func ValidateJwtToken(bt string) (*jwt.Token, error){
	inputToken := RetrieveJwtToken(bt)
	if inputToken == "" {
		return nil, nil // todo, add error here
	}

	token, err := jwt.Parse(inputToken, func(token *jwt.Token) (interface{}, error) {
		// todo add check signing method here
		return []byte(os.Getenv("API_SECRET")), nil
	})

	if err != nil {
		log.Println("failed to validate token")
		return nil, err
	}

	return token, nil
}

func RetrieveUserAuthorization(token *jwt.Token) (*UserAuth, error){
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userUUID, ok := claims["auth_uuid"].(string)
		if !ok {
			err := errors.New("error on type assertion")
			return nil, err
		}

		return &UserAuth{
			UserUUID: userUUID,
			UserID:   int(claims["user_id"].(float64)),
		}, nil
	}
	return nil, errors.New("failed to retrieve user authorization")
}