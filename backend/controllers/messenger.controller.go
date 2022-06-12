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

func GetListActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var page string
	if vars["pid"] == "" {
		page = "1"
	} else {
		page = vars["pid"]
	}
	var profile []byte
	var member view.DatabaseMember
	var check_user bool

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
		check_user = model.CheckUserIsJoinGroup(vars["id"], member.Email)
	}

	if check_user {
		temp, success := model.GetListRoomActivity(vars["id"], page, member.Email)
		if !success {
			temp = []byte("[]")
		}
		resp, _ := json.Marshal(view.HttpSuccessMessage{
			Status:  200,
			Message: temp,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
		return
	} else {
		resp, _ := json.Marshal(view.HttpErrorMessage{
			Status:  403,
			Message: "your are not in the member group list",
		})
		w.WriteHeader(403)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}
