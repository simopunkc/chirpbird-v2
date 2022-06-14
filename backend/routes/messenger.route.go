package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/simopunkc/chirpbird-v2/controllers"
)

func MessengerRoute(app fiber.Router) {
	app.Get("/:room/page:pages", controller.GetMemberChat)
}
