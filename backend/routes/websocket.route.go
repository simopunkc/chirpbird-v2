package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	controller "github.com/simopunkc/chirpbird-v2/controllers"
)

func WebsocketRoute(app fiber.Router) {
	app.Get("/", websocket.New(controller.HandleWebsocket))
}
