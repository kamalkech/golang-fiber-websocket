package main

import (
	"dailycode/learn-fiber/user"
	"dailycode/learn-fiber/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	app := fiber.New()

	// Enable cors.
	app.Use(cors.New())

	// Enable websocket.
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Routes
	user.Routes(app)

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed")) // true
		log.Println(c.Query("v"))        // 1.0

		// Store the WebSocket connection
		ws.WsConnectionsMutex.Lock()
		ws.WsConnections[c] = struct{}{}
		ws.WsConnectionsMutex.Unlock()

		// Send initial user list
		// sendUsers(c)

		// Keep the connection open
		for {
			_, msg, err := c.ReadMessage()
			log.Println("MSG:", msg)

			if err != nil {
				log.Println("read:", err)
				break
			}
		}

		// Remove the WebSocket connection on close
		ws.WsConnectionsMutex.Lock()
		delete(ws.WsConnections, c)
		ws.WsConnectionsMutex.Unlock()
	}))

	log.Fatal(app.Listen(":3333"))

}
