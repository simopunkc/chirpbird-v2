package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func CheckHeaderAccessToken(c *fiber.Ctx) error {
	acc_token := string(c.Request().Header.Peek(os.Getenv("COOKIE_ACCESS_TOKEN")))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		c.Status(401)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		c.Locals("token", string(module.DecryptJWT(acc_token)))
		return c.Next()
	}
}

func CheckDecodeAccessToken(c *fiber.Ctx) error {
	token := c.Locals("token").(string)
	if token != "" {
		temp1, _ := base64.StdEncoding.DecodeString(token)
		c.Locals("decode", temp1)
		return c.Next()
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func CheckValidateAccessToken(c *fiber.Ctx) error {
	var member view.DatabaseMember
	json.Unmarshal(c.Locals("decode").([]byte), &member)
	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		c.Status(403)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		c.Locals("profile", member)
		return c.Next()
	}
}
