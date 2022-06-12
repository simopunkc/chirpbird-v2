package modules

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func fetchGet(url string, id_token string) []byte {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Authorization", "Bearer "+id_token)
	client := &http.Client{}
	resp, _ := client.Do(request)
	b, _ := ioutil.ReadAll(resp.Body)
	time.Sleep(1 * time.Second)
	return b
}

func fetchPost(url string, data []byte) []byte {
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, _ := client.Do(request)
	b, _ := ioutil.ReadAll(resp.Body)
	time.Sleep(1 * time.Second)
	return b
}

func GetGoogleAuthURL(anticsrf string) string {
	return "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + os.Getenv("OAUTH_CLIENT_ID") + "&response_type=code&scope=https://www.googleapis.com/auth/userinfo.profile+https://www.googleapis.com/auth/userinfo.email&redirect_uri=" + os.Getenv("FRONTEND_PROTOCOL") + os.Getenv("FRONTEND_HOST") + os.Getenv("OAUTH_REDIRECT_PATH") + "&state=" + anticsrf + "&access_type=offline&prompt=consent"
}

func RequestGoogleAccessToken(values []byte) []byte {
	const url = "https://oauth2.googleapis.com/token"
	return fetchPost(url, values)
}

func RequestGoogleUserProfile(access_token string, id_token string) string {
	googleUser := fetchGet("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token="+access_token, id_token)
	encodedAccessToken := EncryptJWT(googleUser)
	return encodedAccessToken
}
