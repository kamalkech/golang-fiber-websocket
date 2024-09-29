package main

import (
	"dailycode/learn-fiber/auth"
	"dailycode/learn-fiber/comment"
	"dailycode/learn-fiber/user"
	"dailycode/learn-fiber/ws"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	app := fiber.New()

	// Set max connections
	app.Server().MaxConnsPerIP = 1

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
	comment.Routes(app)

	// Login route
	app.Post("/login", auth.Login)

	// JWT Middleware
	var jwtMiddleware = jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte("secret"),
		},
	})
	app.Use(jwtMiddleware)

	// Unauthenticated route
	app.Get("/accessible", auth.Accessible)
	app.Get("/restricted", jwtMiddleware, auth.Restricted)
	app.Get("/admin", auth.RoleGuard("admin"), auth.AdminOnly)
	app.Get("/user", auth.RoleGuard("user"), auth.UserOnly)
	app.Get("/any", auth.Restricted)

	// Websocket
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed")) // true
		log.Println(c.Query("v"))        // 1.0

		// Store the WebSocket connection
		ws.WsConnectionsMutex.Lock()
		ws.WsConnections[c] = struct{}{}
		ws.WsConnectionsMutex.Unlock()

		// Send initial user list
		comment.SendComments(c)
		user.SendUsers(c)

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

	// Start server
	log.Fatal(app.Listen(":3333"))
}
