package controllers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	model "github.com/simopunkc/chirpbird-v2/models"
	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

func GetSingleActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
	var roomActivity view.DatabaseRoomActivity

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
		rawData := model.RedisReadRoomActivity(vars["id"])
		temp3, _ := base64.StdEncoding.DecodeString(rawData)
		json.Unmarshal(temp3, &roomActivity)
	}

	if roomActivity.Id_room == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  404,
			Message: "Chat not found",
		})
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func DeleteChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var profile []byte
	var member view.DatabaseMember
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
		check_action = model.DeleteRoomActivity(vars["id"], member.Email)
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
			Message: "Failed delete chat",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}
