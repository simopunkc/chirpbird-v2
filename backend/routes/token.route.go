package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/simopunkc/chirpbird-v2/controllers"
)

func TokenRoute(app fiber.Router) {
	app.Get("/csrf", controller.GetCSRF)
}
