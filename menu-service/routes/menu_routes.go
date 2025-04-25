package routes

import (
	"github.com/gofiber/fiber/v2"
	"monolith/menu-service/handlers"
)

func SetupMenuRoutes(app *fiber.App) {
	menu := app.Group("/api/menu")

	// 📦 Склад
	menu.Get("/inventory", handlers.GetInventoryItems)
	menu.Post("/inventory", handlers.CreateInventoryItem)
	menu.Put("/inventory/:id", handlers.UpdateInventoryItem)
	menu.Delete("/inventory/:id", handlers.DeleteInventoryItem)

	// 📊 Калькуляция блюда
	menu.Get("/calculation/:menuItemId", handlers.GetCalculationByMenuItemID)
	menu.Post("/calculation", handlers.CreateCalculationForDish)

	// ✅ Правильный маршрут внутри группы menu
	menu.Get("/published-with-category", handlers.GetPublishedMenuItemsWithCategory)

	// 📥 Опубликованные блюда (видят пользователи)
	menu.Get("/published", handlers.GetPublishedMenuItems)

	// 🔄 Публикация / снятие публикации
	menu.Post("/:id/publish", handlers.PublishMenuItem)

	// 📥 Блюда (админка)
	menu.Get("/with-category", handlers.GetAllMenuItemsWithCategory)
	menu.Get("/", handlers.GetAllMenuItems)
	menu.Post("/", handlers.CreateMenuItem)
	menu.Put("/:id", handlers.UpdateMenuItem)
	menu.Delete("/:id", handlers.DeleteMenuItem)
	menu.Get("/:id", handlers.GetMenuItemByID)
}






