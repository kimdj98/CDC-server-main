package router

import (
	"github.com/gofiber/fiber/v2"
)

// routing to check if server is on
func HelloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello world!")
}
