package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	controller "github.com/simopunkc/chirpbird-v2/controllers"
	module "github.com/simopunkc/chirpbird-v2/modules"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	r := mux.NewRouter()
	r.HandleFunc("/oauth/google/url", controller.GetEndpointURL).Methods("GET")
	r.HandleFunc("/oauth/google/verify", controller.VerifyLogin).Methods("POST")
	r.HandleFunc("/oauth/google/refresh", controller.RefreshLogin).Methods("POST")
	r.HandleFunc("/oauth/google/profile", controller.GetGoogleProfile).Methods("GET")
	r.HandleFunc("/oauth/logout", controller.LogOut).Methods("POST")
	r.HandleFunc("/token/csrf", controller.GetCSRF).Methods("GET")
	r.HandleFunc("/member/room/page{pid:[0-9]+}", controller.GetListRoom).Methods("GET")
	r.HandleFunc("/member/room/create", controller.CreateRoom).Methods("POST")
	r.HandleFunc("/member/room/join", controller.JoinRoom).Methods("PUT")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}", controller.GetSingleRoom).Methods("GET")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/exit", controller.ExitRoom).Methods("PUT")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/rename", controller.RenameRoom).Methods("PUT")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/memberToModerator", controller.MemberToModerator).Methods("PUT")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/ModeratorToMember", controller.ModeratorToMember).Methods("PUT")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/kickMember", controller.KickMember).Methods("PUT")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/addMember", controller.AddMember).Methods("PUT")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/newChat", controller.CreateChat).Methods("POST")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/enableNotif", controller.EnableNotif).Methods("PUT")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/disableNotif", controller.DisableNotif).Methods("PUT")
	r.HandleFunc("/room/{id:[A-Za-z0-9]+}/deleteRoom", controller.DeleteRoom).Methods("DELETE")
	r.HandleFunc("/messenger/{id:[A-Za-z0-9]+}/page{pid:[0-9]+}", controller.GetListActivity).Methods("GET")
	r.HandleFunc("/activity/{id:[A-Za-z0-9]+}", controller.GetSingleActivity).Methods("GET")
	r.HandleFunc("/activity/{id:[A-Za-z0-9]+}/deleteChat", controller.DeleteChat).Methods("DELETE")
	flag.Parse()
	hub := module.NewHub()
	go hub.Run()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		module.ServeWs(hub, w, r)
	})

	fmt.Println("server running at port 9001")
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{
			os.Getenv("FRONTEND_HOST"),
		},
		AllowCredentials: true,
	}).Handler(r)

	backend_host := ":" + os.Getenv("BACKEND_PORT")
	http.ListenAndServe(backend_host, handler)
}
