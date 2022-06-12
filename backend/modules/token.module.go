package modules

import (
	"encoding/json"

	view "github.com/simopunkc/chirpbird-v2/views"
)

func GenerateOauthCsrfToken(miliseconds int64) view.RawData {
	unique := GenerateUniqueID("T")
	timestamp := GetCurrentTimestamp() + miliseconds
	csrftoken := view.Csrftoken{
		Random:  unique,
		Expired: timestamp,
	}
	myJson, _ := json.Marshal(csrftoken)
	encrypt_token := EncryptJWT(myJson)
	rawData := view.RawData{
		Plaintext:  unique,
		Ciphertext: encrypt_token,
	}
	return rawData
}

func ValidateOauthCsrfToken(token string) bool {
	dump := DecryptJWT(token)
	if len(dump) > 0 {
		var temp view.Csrftoken
		err := json.Unmarshal(dump, &temp)
		if err != nil {
			return false
		}
		if temp.Expired >= GetCurrentTimestamp() {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}

func ValidateRefreshToken(expired int64) bool {
	if expired >= GetCurrentTimestamp() {
		return false
	} else {
		return true
	}
}
