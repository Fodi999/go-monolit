package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware — логирует каждый запрос: метод, путь, статус, время выполнения
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next() // продолжаем выполнение

		duration := time.Since(start)
		status := c.Response().StatusCode()
		log.Printf("➡️ %s %s | %d | %s", c.Method(), c.OriginalURL(), status, duration)

		if err != nil {
			log.Printf("❌ Ошибка в обработке: %v", err)
		}

		return err
	}
}
