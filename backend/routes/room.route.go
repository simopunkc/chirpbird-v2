package routes

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/simopunkc/chirpbird-v2/controllers"
)

func RoomRoute(app fiber.Router) {
	app.Get("/:room", controller.GetSingleRoom)
	app.Put("/:room/addMember", controller.PutAddMember)
	app.Put("/:room/exit", controller.PutExitRoom)
	app.Put("/:room/rename", controller.PutRenameRoom)
	app.Put("/:room/memberToModerator", controller.PutMemberToModerator)
	app.Put("/:room/ModeratorToMember", controller.PutModeratorToMember)
	app.Put("/:room/kickMember", controller.PutKickMember)
	app.Post("/:room/newChat", controller.PostNewChat)
	app.Put("/:room/enableNotif", controller.PutEnableNotif)
	app.Put("/:room/disableNotif", controller.PutDisableNotif)
	app.Delete("/:room/deleteRoom", controller.DeleteRoom)
}
