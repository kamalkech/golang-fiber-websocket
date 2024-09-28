package main

import (
	"dailycode/learn-fiber/user"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"sync"
	"time"
)

var (
	users              = make(map[int]*user.User)
	usersMutex         sync.RWMutex
	wsConnections      = make(map[*websocket.Conn]struct{})
	wsConnectionsMutex sync.RWMutex
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

	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(user.User)

		if err := c.BodyParser(user); err != nil {
			return err
		}

		usersMutex.Lock()
		users[user.ID] = user
		usersMutex.Unlock()

		time.Sleep(2 * time.Second)

		triggerFetchUsers()

		return c.JSON(user)
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed")) // true
		log.Println(c.Query("v"))        // 1.0

		// Store the WebSocket connection
		wsConnectionsMutex.Lock()
		wsConnections[c] = struct{}{}
		wsConnectionsMutex.Unlock()

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
		wsConnectionsMutex.Lock()
		delete(wsConnections, c)
		wsConnectionsMutex.Unlock()

	}))

	log.Fatal(app.Listen(":3333"))

}

func triggerFetchUsers() {
	// Trigger all connected WebSocket clients to fetch the updated user list
	wsConnectionsMutex.RLock()
	defer wsConnectionsMutex.RUnlock()

	for conn := range wsConnections {
		sendUsers(conn)
	}
}

func sendUsers(c *websocket.Conn) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	userList := make([]*user.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	if err := c.WriteJSON(userList); err != nil {
		log.Println("write:", err)
	}
}
