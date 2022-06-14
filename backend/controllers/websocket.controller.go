package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"

	model "github.com/simopunkc/chirpbird-v2/models"
	module "github.com/simopunkc/chirpbird-v2/modules"
	view "github.com/simopunkc/chirpbird-v2/views"
)

var ctx = context.Background()

type client struct {
	user view.DatabaseMember
}

var clients = make(map[*websocket.Conn]client)
var register = make(chan *websocket.Conn)
var broadcast = make(chan []byte)
var unregister = make(chan *websocket.Conn)

func RunHub() {
	for {
		select {
		case connection := <-register:
			var profile view.DatabaseMember
			decodeToken := module.DecryptJWT(connection.Cookies(os.Getenv("COOKIE_ACCESS_TOKEN")))
			if len(decodeToken) > 0 {
				temp, _ := base64.StdEncoding.DecodeString(string(decodeToken))
				json.Unmarshal(temp, &profile)
			}
			clients[connection] = client{
				user: profile,
			}
			log.Println(profile.Email, "connection registered")

		case message := <-broadcast:
			BroadcastToOtherServerUsingPubSub(message)

		case connection := <-unregister:
			delete(clients, connection)
			log.Println("connection unregistered")
		}
	}
}

func BroadcastToOtherServerUsingPubSub(message []byte) {
	rdb := model.GetRedisModel()
	defer rdb.Close()
	pubsub := rdb.Subscribe(ctx, "chatting")
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = rdb.Publish(ctx, "chatting", message).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func BroadcastReceivedMessageFromPubSub() {
	rdb := model.GetRedisModel()
	defer rdb.Close()
	for {
		pubsub := rdb.Subscribe(ctx, "chatting")
		ch := pubsub.Channel()
		time.AfterFunc(time.Second, func() {
			_ = pubsub.Close()
		})
		temp := ""
		for msg := range ch {
			var inputWebsocket view.BodyRequestWebsocket
			json.Unmarshal([]byte(msg.Payload), &inputWebsocket)
			if temp != inputWebsocket.Id_room_activity {
				temp = inputWebsocket.Id_room_activity
				fmt.Println("Pubsub message received by server ", os.Getenv("BE_ID"), msg.Channel, msg.Payload)
				for connection := range clients {
					if clients[connection].user.Verified_email {
						check := model.CheckUserIsJoinGroup(inputWebsocket.Id_room, clients[connection].user.Email)
						if check {
							if err := connection.WriteMessage(1, []byte(msg.Payload)); err != nil {
								log.Println("write error:", err)
								connection.WriteMessage(websocket.CloseMessage, []byte{})
								connection.Close()
								delete(clients, connection)
							}
						}
					}
				}
			}
		}
	}
}

func HandleWebsocket(c *websocket.Conn) {
	defer func() {
		unregister <- c
		c.Close()
	}()

	register <- c

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}
			return
		}
		if messageType == websocket.TextMessage {
			broadcast <- message
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}
