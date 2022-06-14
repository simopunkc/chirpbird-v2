package controllers

import (
	"encoding/base64"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	model "github.com/simopunkc/chirpbird-v2/models"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func GetRoomActivity(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var roomActivity view.DatabaseRoomActivity

	rawData := model.RedisReadRoomActivity(c.Params("activity"))
	temp3, _ := base64.StdEncoding.DecodeString(rawData)
	json.Unmarshal(temp3, &roomActivity)

	if roomActivity.Id_room == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Chat not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: map[string]interface{}{
				"id_primary":            roomActivity.Id_primary,
				"id_parent":             roomActivity.Id_parent,
				"id_room":               roomActivity.Id_room,
				"id_member_actor":       roomActivity.Id_member_actor,
				"id_member_target":      roomActivity.Id_member_target,
				"type_activity":         roomActivity.Type_activity,
				"message":               roomActivity.Message,
				"date_created":          roomActivity.Date_created,
				"list_id_member_unread": roomActivity.List_id_member_unread,
				"name":                  member.Name,
				"picture":               member.Picture,
				"verified_email":        member.Verified_email,
			},
		})
		c.Status(200)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	}
}

func DeleteRoomActivity(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)

	var check_action bool = model.DeleteRoomActivity(c.Params("activity"), member.Email)

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: true,
		})
		c.Status(200)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed delete chat",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	}
}
