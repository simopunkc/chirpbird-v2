package controllers

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"

	model "github.com/simopunkc/chirpbird-v2/models"
	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func GetMessageActivity(type_activity string, actor_name string, id_target string) string {
	var targetName string
	if id_target == "" {
		targetName = ""
	} else {
		rawData := model.RedisReadMember(id_target)
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		var rawStruct view.DatabaseMember
		json.Unmarshal(temp2, &rawStruct)
		targetName = rawStruct.Name
	}
	return view.GetActivity(type_activity, actor_name, targetName)
}

func GetSingleRoom(c *fiber.Ctx) error {
	var group view.DatabaseRoom

	rawData := model.RedisReadRoom(c.Params("room"))
	temp2, _ := base64.StdEncoding.DecodeString(rawData)
	json.Unmarshal(temp2, &group)

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: group,
		})
		c.Status(200)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func PutAddMember(c *fiber.Ctx) error {
	var check_user_exist bool
	var check_user_joined bool
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestManipulateRoomMember

	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}
	check_user_exist = model.CheckUserIsExist(parseBody.Id_target)

	if check_user_exist && parseBody.Id_target != member.Email {
		check_user_joined = model.CheckUserIsJoinGroup(c.Params("room"), parseBody.Id_target)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "The user to be added is not registered yet",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}

	if check_user_joined || parseBody.Id_target == member.Email {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "The user is already in the group",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		rawData := model.RedisReadRoom(c.Params("room"))
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "add_member"
		log_message := GetMessageActivity(log_type, member.Name, parseBody.Id_target)
		check_action = model.AddMember(c.Params("room"), member.Email, parseBody.Id_target, newID, date_created, group.List_id_member, log_type, log_message)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: map[string]interface{}{
				"id_room_activity": newID,
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
			Message: "Failed add member",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func PutExitRoom(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var group view.DatabaseRoom
	var newID string
	var check_action bool

	rawData := model.RedisReadRoom(c.Params("room"))
	temp2, _ := base64.StdEncoding.DecodeString(rawData)
	json.Unmarshal(temp2, &group)

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "exit_group"
		log_message := GetMessageActivity(log_type, member.Name, "")
		check_action = model.ExitRoom(c.Params("room"), member.Email, newID, date_created, group.List_id_member, log_type, log_message)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: map[string]interface{}{
				"id_room_activity": newID,
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
			Message: "Failed exit from group",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func PutRenameRoom(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestCreateRoom

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
		rawData := model.RedisReadRoom(c.Params("room"))
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "rename_group"
		log_message := GetMessageActivity(log_type, member.Name, "")
		check_action = model.RenameRoom(c.Params("room"), member.Email, parseBody.Name, newID, date_created, group.List_id_member, log_type, log_message)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: map[string]interface{}{
				"id_room_activity": newID,
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
			Message: "Failed to rename group",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func PutMemberToModerator(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var check_user_exist bool
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestManipulateRoomMember

	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}
	check_user_exist = model.CheckUserIsExist(parseBody.Id_target)

	if check_user_exist && parseBody.Id_target != member.Email {
		rawData := model.RedisReadRoom(c.Params("room"))
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "bad request",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "member_to_moderator"
		log_message := GetMessageActivity(log_type, member.Name, parseBody.Id_target)
		check_action = model.MemberBecomeModerator(c.Params("room"), member.Email, parseBody.Id_target, newID, date_created, group.List_id_member, log_type, log_message)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: map[string]interface{}{
				"id_room_activity": newID,
				"id_member_actor":  member.Email,
				"id_room":          group.Id_primary,
			},
		})
		c.Status(200)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed edit member",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json")
		return c.Send(resp)
	}
}

func PutModeratorToMember(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var check_user_exist bool
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestManipulateRoomMember

	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}
	check_user_exist = model.CheckUserIsExist(parseBody.Id_target)

	if check_user_exist && parseBody.Id_target != member.Email {
		rawData := model.RedisReadRoom(c.Params("room"))
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "bad request",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "moderator_to_member"
		log_message := GetMessageActivity(log_type, member.Name, parseBody.Id_target)
		check_action = model.ModeratorBecomeMember(c.Params("room"), member.Email, parseBody.Id_target, newID, date_created, group.List_id_member, log_type, log_message)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: map[string]interface{}{
				"id_room_activity": newID,
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
			Message: "Failed edit member",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func PutKickMember(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var check_user_exist bool
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestManipulateRoomMember

	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}
	check_user_exist = model.CheckUserIsExist(parseBody.Id_target)

	if check_user_exist && parseBody.Id_target != member.Email {
		rawData := model.RedisReadRoom(c.Params("room"))
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "bad request",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "kick_member"
		log_message := GetMessageActivity(log_type, member.Name, parseBody.Id_target)
		check_action = model.KickMember(c.Params("room"), member.Email, parseBody.Id_target, newID, date_created, group.List_id_member, log_type, log_message)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 200,
			Message: map[string]interface{}{
				"id_room_activity": newID,
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
			Message: "Failed to kick member",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func PostNewChat(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var parseBody view.BodyRequestCreateRoomActivity
	var group view.DatabaseRoom
	var newID string
	var check_publish bool

	err := c.BodyParser(&parseBody)
	if err != nil {
		log.Fatal(err)
	}

	if parseBody.Message == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request message not found",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		rawData := model.RedisReadRoom(c.Params("room"))
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Group not found",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		id_member_target := ""
		if parseBody.Id_parent != "" {
			var temp2 view.DatabaseRoomActivity
			rawData := model.RedisReadRoomActivity(parseBody.Id_parent)
			temp3, _ := base64.StdEncoding.DecodeString(rawData)
			json.Unmarshal(temp3, &temp2)
			id_member_target = temp2.Id_member_actor
		}
		date_created := module.GenerateIsoDate()
		newID = module.GenerateUniqueID("RA")
		log_type := "new_chat"
		check_publish = model.CreateChat(newID, parseBody.Id_parent, group.Id_primary, member.Email, id_member_target, parseBody.Message, date_created, group.List_id_member, log_type)
	}

	if check_publish {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status: 201,
			Message: map[string]interface{}{
				"id_room_activity": newID,
				"id_member_actor":  member.Email,
				"id_room":          group.Id_primary,
			},
		})
		c.Status(201)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed publish chat",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func PutEnableNotif(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var group view.DatabaseRoom
	var check_action bool

	rawData := model.RedisReadRoom(c.Params("room"))
	temp2, _ := base64.StdEncoding.DecodeString(rawData)
	json.Unmarshal(temp2, &group)

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		check_action = model.EnableRoomNotification(c.Params("room"), member.Email)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: true,
		})
		c.Status(200)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed to enable group notifications",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func PutDisableNotif(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var group view.DatabaseRoom
	var check_action bool

	rawData := model.RedisReadRoom(c.Params("room"))
	temp2, _ := base64.StdEncoding.DecodeString(rawData)
	json.Unmarshal(temp2, &group)

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		check_action = model.DisableRoomNotification(c.Params("room"), member.Email)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: true,
		})
		c.Status(200)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed to disable group notifications",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}

func DeleteRoom(c *fiber.Ctx) error {
	var member view.DatabaseMember = c.Locals("profile").(view.DatabaseMember)
	var group view.DatabaseRoom
	var check_action bool

	rawData := model.RedisReadRoom(c.Params("room"))
	temp2, _ := base64.StdEncoding.DecodeString(rawData)
	json.Unmarshal(temp2, &group)

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		c.Status(404)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		check_action = model.DeleteRoom(c.Params("room"), member.Email)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: true,
		})
		c.Status(200)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed to delete group",
		})
		c.Status(400)
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Send(resp)
	}
}
