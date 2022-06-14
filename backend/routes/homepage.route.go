package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/simopunkc/chirpbird-v2/controllers"
)

func HomepageRoute(app fiber.Router) {
	app.Get("/", controller.GetHomepage)
}
