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

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var profile []byte
	var member view.DatabaseMember
	var parseBody view.BodyRequestCreateRoom
	var id_room string
	var id_log string
	var create_group bool

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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  500,
			Message: "failed create group",
		})
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	}
}

func GetListRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var page string
	if vars["pid"] == "" {
		page = "1"
	} else {
		page = vars["pid"]
	}
	var profile []byte
	var member view.DatabaseMember

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
		temp, success := model.GetListRoom(page, member.Email)
		if !success {
			temp = []byte("[]")
		}
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: temp,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	var profile []byte
	var member view.DatabaseMember
	var parseBody view.BodyRequestJoinRoom
	var check_token bool
	var id_log string
	var join_group bool
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

	if parseBody.Token == "" {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "request token not found",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
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
		w.WriteHeader(404)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  400,
			Message: "Failed joining group",
		})
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}
