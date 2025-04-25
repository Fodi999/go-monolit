package routes

import (
	"github.com/gofiber/fiber/v2"
	"monolith/menu-service/handlers"
)

func SetupCategoryRoutes(app *fiber.App) {
	api := app.Group("/api/categories")

	// Получить все категории
	api.Get("/", handlers.GetAllCategories)

	// Создать категорию
	api.Post("/", handlers.CreateCategory)

	// Обновить категорию по ID
	api.Put("/:id", handlers.UpdateCategory)

	// Удалить категорию по ID
	api.Delete("/:id", handlers.DeleteCategory)
}

