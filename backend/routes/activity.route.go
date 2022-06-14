package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/simopunkc/chirpbird-v2/controllers"
)

func ActivityRoute(app fiber.Router) {
	app.Get("/:activity", controller.GetRoomActivity)
	app.Delete("/:activity/deleteChat", controller.DeleteRoomActivity)
}
