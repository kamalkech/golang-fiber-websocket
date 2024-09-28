package user

import (
	"dailycode/learn-fiber/ws"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var (
	users      = make(map[int]*User)
	usersMutex sync.RWMutex
)

func Routes(app *fiber.App) {
	var UserRoutes = app.Group("/users")

	// Get all users
	UserRoutes.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("list users")
	})

	// Create new user
	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return err
		}

		usersMutex.Lock()
		users[user.ID] = user
		usersMutex.Unlock()

		time.Sleep(2 * time.Second)

		// Trigger to sent user list
		triggerFetchUsers()

		return c.JSON(user)
	})

}

func triggerFetchUsers() {
	// Trigger all connected WebSocket clients to fetch the updated user list
	ws.WsConnectionsMutex.RLock()
	defer ws.WsConnectionsMutex.RUnlock()

	for conn := range ws.WsConnections {
		SendUsers(conn)
	}
}

func SendUsers(c *websocket.Conn) {
	usersMutex.RLock()
	defer usersMutex.RUnlock()

	userList := make([]*User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}

	if err := c.WriteJSON(userList); err != nil {
		log.Println("write:", err)
	}
}
