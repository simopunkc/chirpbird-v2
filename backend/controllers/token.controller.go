package controllers

import (
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"

	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func GetCSRF(c *fiber.Ctx) error {
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
	c.Status(200)
	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.Send(resp)
}
