package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func GetCSRF(w http.ResponseWriter, r *http.Request) {
	csrf := module.GenerateOauthCsrfToken(module.CSRF_TOKEN_TIMEOUT)
	resp, _ := json.Marshal(view.HttpSuccessMessage{
		Status: 200,
		Message: []interface{}{
			map[string]interface{}{
				"cookie": os.Getenv("COOKIE_CSRF"),
				"value":  csrf.Ciphertext,
				"config": map[string]interface{}{
					"maxAge":   module.CSRF_TOKEN_TIMEOUT,
					"httpOnly": false,
					"secure":   false,
					"domain":   os.Getenv("FRONTEND_DOMAIN"),
				},
			},
		},
	})
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
