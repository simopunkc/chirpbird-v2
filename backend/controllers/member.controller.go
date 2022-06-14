package controllers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	model "github.com/simopunkc/chirpbird-v2/models"
	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func GetMemberRoom(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)

	var page string
	if c.Params("pages") == "" {
		page = "1"
	} else {
		page = c.Params("pages")
	}

	temp, success := model.GetListRoom(page, member.Email)
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
}

func PostCreateRoom(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var parseBody view.BodyRequestCreateRoom
	var id_room string
	var id_log string
	var create_group bool

	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}

	if parseBody.Name == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request name not found",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		id_room = module.GenerateUniqueID("G")
		link_group := module.GenerateUniqueID("LJ")
		date_created := module.GenerateIsoDate()
		id_log = module.GenerateUniqueID("RA")
		log_type := "group_created"
		log_message := GetMessageActivity(log_type, member.Name, "")
		create_group = model.CreateRoom(id_room, member.Email, parseBody.Name, []string{member.Email}, []string{member.Email}, []string{}, []string{}, date_created, date_created, link_group, id_log, log_message, log_type)
	}

	if create_group {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 201,
			Message: map[string]interface{}{
				"id_room_activity": id_log,
				"id_member_actor":  member.Email,
				"id_room":          id_room,
			},
		})
		c.Status(200)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  500,
			Message: "failed create group",
		})
		c.Status(500)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func PutJoinRoom(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var parseBody view.BodyRequestJoinRoom
	var check_token bool
	var id_log string
	var join_group bool
	var group view.DatabaseRoom

	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}

	if parseBody.Token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request token not found",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		temp, check := model.CheckTokenGroupExist(parseBody.Token)
		if check {
			check_token = true
			final, _ := json.Marshal(temp)
			json.Unmarshal(final, &group)
		} else {
			check_token = false
		}
	}

	if check_token {
		id_log = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "join_group"
		log_message := GetMessageActivity(log_type, member.Name, "")
		join_group = model.JoinRoom(group.Id_primary, member.Email, id_log, date_created, group.List_id_member, log_message, log_type)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}

	if join_group {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: map[string]interface{}{
				"id_room_activity": id_log,
				"id_member_actor":  member.Email,
				"id_room":          group.Id_primary,
			},
		})
		c.Status(200)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed joining group",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}
