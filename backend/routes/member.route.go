package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/simopunkc/chirpbird-v2/controllers"
)

func MemberRoute(app fiber.Router) {
	app.Get("/room/page:pages", controller.GetMemberRoom)
	app.Post("/room/create", controller.PostCreateRoom)
	app.Put("/room/join", controller.PutJoinRoom)
}
