package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/simopunkc/chirpbird-v2/controllers"
	middleware "github.com/simopunkc/chirpbird-v2/middlewares"
)

func OauthGoogleRoute(app fiber.Router) {
	app.Get("/google/url", controller.GetLoginPage)
	app.Post("/google/verify", controller.GetVerifyLogin)
	app.Post("/google/refresh", controller.GetRefreshLogin)
	app.Post("/logout", controller.GetLogOut)
	app.Get("/google/profile", middleware.CheckHeaderAccessToken, middleware.CheckDecodeAccessToken, middleware.CheckValidateAccessToken, controller.GetProfileLogin)
}
