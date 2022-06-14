package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebsocketUpgrader(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}
