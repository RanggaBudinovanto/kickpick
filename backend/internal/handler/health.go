package handler

import "github.com/gofiber/fiber/v2"

func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

func NotImplemented(c *fiber.Ctx) error {
	return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
		"error": "database belum terhubung",
	})
}
