package controllers

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	model "github.com/simopunkc/chirpbird-v2/models"
	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func GetEndpointURL(w http.ResponseWriter, r *http.Request) {
	csrf := module.GenerateOauthCsrfToken(module.ACC_TOKEN_TIMEOUT)
	link := module.GetGoogleAuthURL(csrf.Plaintext)
	resp, _ := json.Marshal(view.HttpSuccessMessage{
		Status: 200,
		Message: []interface{}{
			map[string]interface{}{
				"url": link,
				"state": map[string]interface{}{
					"cookie": os.Getenv("COOKIE_CSRF"),
					"value":  csrf.Ciphertext,
					"config": map[string]interface{}{
						"maxAge":   module.ACC_TOKEN_TIMEOUT,
						"httpOnly": false,
						"secure":   false,
						"domain":   os.Getenv("FRONTEND_DOMAIN"),
					},
				},
			},
		},
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func VerifyLogin(w http.ResponseWriter, r *http.Request) {
	var anticsrf string
	var csrf_token view.Csrftoken
	var validCsrf bool
	var xsrfToken []byte
	var userGoogle []byte
	var isExpired bool
	var encryptRefreshToken string
	var encryptJWTGoogleProfile string
	var member view.DatabaseMember

	b, _ := ioutil.ReadAll(r.Body)
	var parseBody view.BodyRequestVerifyLogin
	err := json.Unmarshal(b, &parseBody)
	if err != nil {
		log.Fatal(err)
	}
	xsrf := r.Header.Get(os.Getenv("COOKIE_CSRF"))
	if parseBody.State == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request state not found",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		anticsrf = parseBody.State
	}

	if xsrf == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "header xsrf not found",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		xsrfToken = module.DecryptJWT(xsrf)
	}

	if len(xsrfToken) > 0 {
		temp2, _ := base64.StdEncoding.DecodeString(string(xsrfToken))
		json.Unmarshal(temp2, &csrf_token)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid xsrf token",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if parseBody.Code == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request code not found",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		if anticsrf == csrf_token.Random {
			validCsrf = true
		} else {
			validCsrf = false
		}
	}

	if validCsrf {
		isExpired = module.ValidateOauthCsrfToken(xsrf)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "xsrf token doesnt match",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if isExpired {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "xsrf token expired, try again",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		paramAuthentikasi := "code=" + parseBody.Code + "&clientId=" + os.Getenv("OAUTH_CLIENT_ID") + "&clientSecret=" + os.Getenv("OAUTH_SECRET") + "&redirectUri=" + os.Getenv("FRONTEND_PROTOCOL") + os.Getenv("FRONTEND_HOST") + os.Getenv("OAUTH_REDIRECT_PATH") + "&grant_type=authorization_code"
		var resp view.BodyPostResponseVerifyLogin
		temp := module.RequestGoogleAccessToken([]byte(paramAuthentikasi))
		json.Unmarshal(temp, &resp)

		encryptJWTGoogleProfile = module.RequestGoogleUserProfile(resp.Access_token, resp.Id_token)

		userGoogle = module.DecryptJWT(encryptJWTGoogleProfile)
		if len(userGoogle) > 0 {
			temp3, _ := base64.StdEncoding.DecodeString(string(userGoogle))
			json.Unmarshal(temp3, &member)

			timestamp := module.GetCurrentTimestamp() + module.REF_TOKEN_TIMEOUT
			temp_token := view.TempRefreshToken{
				Value:   resp.Refresh_token,
				Expired: timestamp,
			}
			temp_refresh_token, _ := json.Marshal(temp_token)
			encryptRefreshToken = module.EncryptJWT(temp_refresh_token)
		} else {
			userGoogle = []byte{}
		}
	}

	if len(userGoogle) == 0 {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid jwt",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "login expired, try again",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		model.SaveMember(member.Email, member.Name, member.Picture, member.Verified_email)

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: []interface{}{
				map[string]interface{}{
					"cookie": os.Getenv("COOKIE_ACCESS_TOKEN"),
					"value":  encryptJWTGoogleProfile,
					"config": map[string]interface{}{
						"maxAge":   module.ACC_TOKEN_TIMEOUT,
						"httpOnly": false,
						"secure":   false,
						"domain":   os.Getenv("FRONTEND_DOMAIN"),
					},
				},
				map[string]interface{}{
					"cookie": os.Getenv("COOKIE_REFRESH_TOKEN"),
					"value":  encryptRefreshToken,
					"config": map[string]interface{}{
						"maxAge":   module.REF_TOKEN_TIMEOUT,
						"httpOnly": false,
						"secure":   false,
						"domain":   os.Getenv("FRONTEND_DOMAIN"),
					},
				},
			},
		})
		w.Write(resp)
	}
}

func RefreshLogin(w http.ResponseWriter, r *http.Request) {
	var exp_xsrf_token bool
	var decodeToken []byte
	var exp_refresh_token bool
	var tempToken view.TempRefreshToken

	xsrf := r.Header.Get(os.Getenv("COOKIE_CSRF"))
	if xsrf == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "header xsrf not found",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		exp_xsrf_token = module.ValidateOauthCsrfToken(xsrf)
	}

	b, _ := ioutil.ReadAll(r.Body)
	var parseBody view.BodyRequestRefreshLogin
	err := json.Unmarshal(b, &parseBody)
	if err != nil {
		log.Fatal(err)
	}

	if parseBody.Ref_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request refresh token header not found",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		decodeToken = module.DecryptJWT(parseBody.Ref_token)
	}

	if len(decodeToken) > 0 {
		temp, _ := base64.StdEncoding.DecodeString(string(decodeToken))
		json.Unmarshal(temp, &tempToken)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid refresh token",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if !exp_xsrf_token {
		exp_refresh_token = module.ValidateRefreshToken(tempToken.Expired)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "xsrf token expired, try again",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if exp_refresh_token {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "refresh token expired, try login again",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		paramRefresh := "client_id=" + os.Getenv("OAUTH_CLIENT_ID") + "&client_secret=" + os.Getenv("OAUTH_SECRET") + "&refresh_token=" + tempToken.Value + "&grant_type=refresh_token"
		var resp view.BodyPostResponseRefreshLogin
		temp2 := module.RequestGoogleAccessToken([]byte(paramRefresh))
		json.Unmarshal(temp2, &resp)

		encryptJWTGoogleProfile := module.RequestGoogleUserProfile(resp.Access_token, resp.Id_token)
		result, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: []interface{}{
				map[string]interface{}{
					"cookie": os.Getenv("COOKIE_ACCESS_TOKEN"),
					"value":  encryptJWTGoogleProfile,
					"config": map[string]interface{}{
						"maxAge":   module.ACC_TOKEN_TIMEOUT,
						"httpOnly": false,
						"secure":   false,
						"domain":   os.Getenv("FRONTEND_DOMAIN"),
					},
				},
			},
		})
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}

func GetGoogleProfile(w http.ResponseWriter, r *http.Request) {
	var profile []byte
	var member view.DatabaseMember

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: member,
		})
		w.Write(resp)
	}
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var parseBody view.BodyRequestRefreshLogin
	err := json.Unmarshal(b, &parseBody)
	if err != nil {
		log.Fatal(err)
	}

	if parseBody.Ref_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request refresh token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: []interface{}{
				map[string]interface{}{
					"cookie": os.Getenv("COOKIE_ACCESS_TOKEN"),
					"value":  "",
					"config": map[string]interface{}{
						"maxAge":   0,
						"httpOnly": false,
						"secure":   false,
						"domain":   os.Getenv("FRONTEND_DOMAIN"),
					},
				},
				map[string]interface{}{
					"cookie": os.Getenv("COOKIE_REFRESH_TOKEN"),
					"value":  "",
					"config": map[string]interface{}{
						"maxAge":   0,
						"httpOnly": false,
						"secure":   false,
						"domain":   os.Getenv("FRONTEND_DOMAIN"),
					},
				},
			},
		})
		w.Write(resp)
	}
}
