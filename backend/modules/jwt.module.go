package modules

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

func EncryptJWT(plaintext []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": plaintext,
	})
	jwt_secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(jwt_secret)
	if err != nil {
		return "error"
	}
	return tokenString
}

func ConvertInterface(rawdata interface{}) []byte {
	final, _ := json.Marshal(rawdata)
	var token string
	json.Unmarshal(final, &token)
	return []byte(token)
}

func DecryptJWT(ciphertext string) []byte {
	token, err := jwt.Parse(ciphertext, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		jwt_secret := []byte(os.Getenv("JWT_SECRET"))
		return jwt_secret, nil
	})
	if err != nil {
		return []byte{}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var hasil []byte
		for key, val := range claims {
			if key == "data" {
				hasil = ConvertInterface(val)
				break
			} else {
				continue
			}
		}
		return hasil
	} else {
		return []byte{}
	}
}
