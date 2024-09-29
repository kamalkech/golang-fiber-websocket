package comment

import (
	"dailycode/learn-fiber/ws"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"log"
	"sync"
	"time"
)

var (
	comments      = make(map[int]*Comment)
	commentsMutex sync.RWMutex
)

func Routes(app *fiber.App) {
	var routes = app.Group("/comments")

	// Get all items
	routes.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("list comments")
	})

	// Create new item
	routes.Post("/", func(c *fiber.Ctx) error {
		comment := new(Comment)

		if err := c.BodyParser(comment); err != nil {
			return err
		}

		commentsMutex.Lock()
		comments[comment.ID] = comment
		commentsMutex.Unlock()

		time.Sleep(2 * time.Second)

		// Trigger to sent list
		triggerFetchItems()

		return c.JSON(comment)
	})

}

func triggerFetchItems() {
	// Trigger all connected WebSocket clients to fetch the updated item list
	ws.WsConnectionsMutex.RLock()
	defer ws.WsConnectionsMutex.RUnlock()

	for conn := range ws.WsConnections {
		SendComments(conn)
	}
}

func SendComments(c *websocket.Conn) {
	commentsMutex.RLock()
	defer commentsMutex.RUnlock()

	commentList := make([]*Comment, 0, len(comments))
	for _, comment := range comments {
		commentList = append(commentList, comment)
	}

	if err := c.WriteJSON(commentList); err != nil {
		log.Println("write:", err)
	}
}
