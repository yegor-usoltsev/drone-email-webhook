package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func main() {
	settings := NewSettingsFromEnv()

	server := NewFiber(settings)

	server.Use(recover.New())

	withTimeout := func(handler fiber.Handler) fiber.Handler {
		return timeout.NewWithContext(handler, settings.ServerHandlerTimeout)
	}

	server.Get("/health", withTimeout(func(c *fiber.Ctx) error {
		return c.SendString("OK")
	}))

	server.ListenGracefully()
}
