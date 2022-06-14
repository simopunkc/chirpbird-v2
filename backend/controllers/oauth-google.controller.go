package controllers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	model "github.com/simopunkc/chirpbird-v2/models"
	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func GetLoginPage(c *fiber.Ctx) error {
	csrf := module.GenerateOauthCsrfToken(module.ACC_TOKEN_TIMEOUT)
	link := module.GetGoogleAuthURL(csrf.Plaintext)
	temp1 := base64.StdEncoding.EncodeToString([]byte(link))
	resp, _ := json.Marshal(view.HttpSuccessMessage{
		Status: 200,
		Message: []interface{}{
			map[string]interface{}{
				"url": temp1,
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
	c.Status(200)
	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.Send(resp)
}

func GetVerifyLogin(c *fiber.Ctx) error {
	var anticsrf string
	var csrf_token view.Csrftoken
	var validCsrf bool
	var xsrfToken []byte
	var userGoogle []byte
	var isExpired bool
	var encryptRefreshToken string
	var respVerify view.BodyPostResponseVerifyLogin
	var encryptJWTGoogleProfile string
	var member view.DatabaseMember

	var parseBody view.BodyRequestVerifyLogin
	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}

	if parseBody.State == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request state not found",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		anticsrf = parseBody.State
	}

	xsrf := string(c.Request().Header.Peek(os.Getenv("COOKIE_CSRF")))
	if xsrf == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "header xsrf not found",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
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
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}

	if parseBody.Code == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request code not found",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
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
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}

	if isExpired {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "xsrf token expired, try again",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		paramAuthentikasi := "code=" + parseBody.Code + "&clientId=" + os.Getenv("OAUTH_CLIENT_ID") + "&clientSecret=" + os.Getenv("OAUTH_SECRET") + "&redirectUri=" + os.Getenv("FRONTEND_PROTOCOL") + os.Getenv("FRONTEND_HOST") + os.Getenv("OAUTH_REDIRECT_PATH") + "&grant_type=authorization_code"
		temp := module.RequestGoogleAccessToken([]byte(paramAuthentikasi))
		json.Unmarshal(temp, &respVerify)
	}

	if respVerify.Access_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "fetch post verify login to google failed",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		encryptJWTGoogleProfile = module.RequestGoogleUserProfile(respVerify.Access_token, respVerify.Id_token)
		userGoogle = module.DecryptJWT(encryptJWTGoogleProfile)
	}

	if len(userGoogle) > 0 {
		temp3, _ := base64.StdEncoding.DecodeString(string(userGoogle))
		json.Unmarshal(temp3, &member)

		timestamp := module.GetCurrentTimestamp() + module.REF_TOKEN_TIMEOUT
		temp_token := view.TempRefreshToken{
			Value:   respVerify.Refresh_token,
			Expired: timestamp,
		}
		temp_refresh_token, _ := json.Marshal(temp_token)
		encryptRefreshToken = module.EncryptJWT(temp_refresh_token)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "failed decrypt jwt google profile",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "login expired, try again",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		model.SaveMember(member.Email, member.Name, member.Picture, member.Verified_email)

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
		c.Status(200)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func GetRefreshLogin(c *fiber.Ctx) error {
	var exp_xsrf_token bool
	var decodeToken []byte
	var exp_refresh_token bool
	var respRefresh view.BodyPostResponseRefreshLogin
	var tempToken view.TempRefreshToken

	xsrf := string(c.Request().Header.Peek(os.Getenv("COOKIE_CSRF")))
	if xsrf == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "header xsrf not found",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		exp_xsrf_token = module.ValidateOauthCsrfToken(xsrf)
	}

	var parseBody view.BodyRequestRefreshLogin
	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}

	if parseBody.Ref_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request refresh token header not found",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
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
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}

	if !exp_xsrf_token {
		exp_refresh_token = module.ValidateRefreshToken(tempToken.Expired)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "xsrf token expired, try again",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}

	if exp_refresh_token {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "refresh token expired, try login again",
		})
		c.Status(403)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		paramRefresh := "client_id=" + os.Getenv("OAUTH_CLIENT_ID") + "&client_secret=" + os.Getenv("OAUTH_SECRET") + "&refresh_token=" + tempToken.Value + "&grant_type=refresh_token"
		temp2 := module.RequestGoogleAccessToken([]byte(paramRefresh))
		json.Unmarshal(temp2, &respRefresh)
	}

	if respRefresh.Access_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "failed fetch refreh token to google",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		encryptJWTGoogleProfile := module.RequestGoogleUserProfile(respRefresh.Access_token, respRefresh.Id_token)
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
			},
		})
		c.Status(200)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func GetLogOut(c *fiber.Ctx) error {
	var parseBody view.BodyRequestRefreshLogin
	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}

	if parseBody.Ref_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request refresh token header not found",
		})
		c.Status(401)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
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
		c.Status(200)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func GetProfileLogin(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)

	resp, _ := json.Marshal(view.HttpSuccessMessage{
		Status:  200,
		Message: member,
	})
	c.Status(200)
	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.Send(resp)
}
