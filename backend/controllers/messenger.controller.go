package controllers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"

	model "github.com/simopunkc/chirpbird-v2/models"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func GetMemberChat(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)

	var page string
	if c.Params("pages") == "" {
		page = "1"
	} else {
		page = c.Params("pages")
	}

	var check_user bool = model.CheckUserIsJoinGroup(c.Params("room"), member.Email)

	if check_user {
		temp, success := model.GetListRoomActivity(c.Params("room"), page, member.Email)
		if !success {
			temp = []byte("[]")
		}
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: temp,
		})
		c.Status(200)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your are not in the member group list",
		})
		c.Status(403)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	}
}
