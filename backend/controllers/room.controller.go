package controllers

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	model "github.com/simopunkc/chirpbird-v2/models"
	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func CreateChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var parseBody view.BodyRequestCreateRoomActivity
	var group view.DatabaseRoom
	var newID string
	var check_publish bool

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &parseBody)
		if err != nil {
			log.Fatal(err)
		}
	}

	if parseBody.Message == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request message not found",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Group not found",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed publish chat",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func GetSingleRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var group view.DatabaseRoom

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: group,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

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

func ExitRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var group view.DatabaseRoom
	var newID string
	var check_action bool

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "exit_group"
		log_message := GetMessageActivity(log_type, member.Name, "")
		check_action = model.ExitRoom(vars["id"], member.Email, newID, date_created, group.List_id_member, log_type, log_message)
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed exit from group",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func RenameRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestCreateRoom

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &parseBody)
		if err != nil {
			log.Fatal(err)
		}
	}

	if parseBody.Name == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request name not found",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "rename_group"
		log_message := GetMessageActivity(log_type, member.Name, "")
		check_action = model.RenameRoom(vars["id"], member.Email, parseBody.Name, newID, date_created, group.List_id_member, log_type, log_message)
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed to rename group",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func KickMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var check_user_exist bool
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestManipulateRoomMember

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &parseBody)
		if err != nil {
			log.Fatal(err)
		}
		check_user_exist = model.CheckUserIsExist(parseBody.Id_target)
	}

	if check_user_exist && parseBody.Id_target != member.Email {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "bad request",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "kick_member"
		log_message := GetMessageActivity(log_type, member.Name, parseBody.Id_target)
		check_action = model.KickMember(vars["id"], member.Email, parseBody.Id_target, newID, date_created, group.List_id_member, log_type, log_message)
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed to kick member",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func AddMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var check_user_exist bool
	var check_user_joined bool
	var member view.DatabaseMember
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestManipulateRoomMember

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &parseBody)
		if err != nil {
			log.Fatal(err)
		}
		check_user_exist = model.CheckUserIsExist(parseBody.Id_target)
	}

	if check_user_exist && parseBody.Id_target != member.Email {
		check_user_joined = model.CheckUserIsJoinGroup(vars["id"], parseBody.Id_target)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "The user to be added is not registered yet",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if check_user_joined || parseBody.Id_target == member.Email {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "The user is already in the group",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "add_member"
		log_message := GetMessageActivity(log_type, member.Name, parseBody.Id_target)
		check_action = model.AddMember(vars["id"], member.Email, parseBody.Id_target, newID, date_created, group.List_id_member, log_type, log_message)
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed add member",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func MemberToModerator(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var check_user_exist bool
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestManipulateRoomMember

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &parseBody)
		if err != nil {
			log.Fatal(err)
		}
		check_user_exist = model.CheckUserIsExist(parseBody.Id_target)
	}

	if check_user_exist && parseBody.Id_target != member.Email {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "bad request",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "member_to_moderator"
		log_message := GetMessageActivity(log_type, member.Name, parseBody.Id_target)
		check_action = model.MemberBecomeModerator(vars["id"], member.Email, parseBody.Id_target, newID, date_created, group.List_id_member, log_type, log_message)
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed edit member",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func ModeratorToMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var check_user_exist bool
	var group view.DatabaseRoom
	var newID string
	var check_action bool
	var parseBody view.BodyRequestManipulateRoomMember

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		b, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(b, &parseBody)
		if err != nil {
			log.Fatal(err)
		}
		check_user_exist = model.CheckUserIsExist(parseBody.Id_target)
	}

	if check_user_exist && parseBody.Id_target != member.Email {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "bad request",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		newID = module.GenerateUniqueID("RA")
		date_created := module.GenerateIsoDate()
		log_type := "moderator_to_member"
		log_message := GetMessageActivity(log_type, member.Name, parseBody.Id_target)
		check_action = model.ModeratorBecomeMember(vars["id"], member.Email, parseBody.Id_target, newID, date_created, group.List_id_member, log_type, log_message)
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed edit member",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func EnableNotif(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var group view.DatabaseRoom
	var check_action bool

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		check_action = model.EnableRoomNotification(vars["id"], member.Email)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed to enable group notifications",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func DisableNotif(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var group view.DatabaseRoom
	var check_action bool

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		check_action = model.DisableRoomNotification(vars["id"], member.Email)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed to disable group notifications",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var group view.DatabaseRoom
	var check_action bool

	acc_token := r.Header.Get(os.Getenv("COOKIE_ACCESS_TOKEN"))
	if acc_token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  401,
			Message: "request access token header not found",
		})
		w.WriteHeader(401)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		profile = module.DecryptJWT(acc_token)
	}

	if len(profile) > 0 {
		temp1, _ := base64.StdEncoding.DecodeString(string(profile))
		json.Unmarshal(temp1, &member)
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "invalid access token header",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}

	if member.Email == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your email is not registered in our database",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		rawData := model.RedisReadRoom(vars["id"])
		temp2, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp2, &group)
	}

	if group.Id_member_creator == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Group not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		check_action = model.DeleteRoom(vars["id"], member.Email)
	}

	if check_action {
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed to delete group",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}
