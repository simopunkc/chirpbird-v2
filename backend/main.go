package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"github.com/joho/godotenv"

	controller "github.com/simopunkc/chirpbird-v2/controllers"
	middleware "github.com/simopunkc/chirpbird-v2/middlewares"
	router "github.com/simopunkc/chirpbird-v2/routes"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	backend_host := ":" + os.Getenv("BACKEND_PORT")

	config := fiber.Config{
		ServerHeader:  "Pasuh",
		StrictRouting: true,
		CaseSensitive: true,
	}

	app := fiber.New(config)

	homepageRoute := app.Group("", middleware.SetSecurityHeader)
	router.HomepageRoute(homepageRoute)

	tokenRoute := app.Group("/token", middleware.SetSecurityHeader)
	router.TokenRoute(tokenRoute)

	oauthGoogleRoute := app.Group("/oauth", middleware.SetSecurityHeader)
	router.OauthGoogleRoute(oauthGoogleRoute)

	memberRoute := app.Group("/member", middleware.SetSecurityHeader, middleware.CheckHeaderAccessToken, middleware.CheckDecodeAccessToken, middleware.CheckValidateAccessToken)
	router.MemberRoute(memberRoute)

	roomRoute := app.Group("/room", middleware.SetSecurityHeader, middleware.CheckHeaderAccessToken, middleware.CheckDecodeAccessToken, middleware.CheckValidateAccessToken)
	router.RoomRoute(roomRoute)

	activityRoute := app.Group("/activity", middleware.SetSecurityHeader, middleware.CheckHeaderAccessToken, middleware.CheckDecodeAccessToken, middleware.CheckValidateAccessToken)
	router.ActivityRoute(activityRoute)

	messengerRoute := app.Group("/messenger", middleware.SetSecurityHeader, middleware.CheckHeaderAccessToken, middleware.CheckDecodeAccessToken, middleware.CheckValidateAccessToken)
	router.MessengerRoute(messengerRoute)

	go controller.RunHub()
	go controller.BroadcastReceivedMessageFromPubSub()

	websocketRoute := app.Group("/ws", middleware.WebsocketUpgrader)
	router.WebsocketRoute(websocketRoute)

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join([]string{
			os.Getenv("FRONTEND_PROTOCOL") + os.Getenv("FRONTEND_HOST"),
		}, ", "),
	}))

	app.Use(limiter.New(limiter.Config{
		Max:               30,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	fmt.Println("server running")

	addr := flag.String("addr", backend_host, "http service address")
	flag.Parse()
	log.Fatal(app.Listen(*addr))
}
