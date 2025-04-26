package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	// Auth
	authDB "monolith/auth_service/database"
	authHandlers "monolith/auth_service/handlers"
	authMiddleware "monolith/auth_service/middleware"

	// Menu
	menuDB "monolith/menu-service/database"
	menuHandlers "monolith/menu-service/handlers"
	menuMiddleware "monolith/menu-service/middleware"
)

func main() {
	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env файл не найден — используем переменные окружения")
	}

	// Инициализация баз данных
	authDB.Init()
	menuDB.InitPostgres()

	// Инициализация Fiber с глобальной обработкой ошибок
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Глобальные middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(logger.New())
	app.Use(menuMiddleware.LoggerMiddleware())

	// ✅ Базовый маршрут "/" для проверки работоспособности сервера
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("✅ Backend is up and running")
	})

	// ✅ Обработка favicon.ico (без ошибки)
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	// === AUTH ===
	app.Post("/api/register", authHandlers.HandleRegister)
	app.Post("/api/login", authHandlers.HandleLogin)
	app.Get("/api/check-email", authHandlers.HandleCheckEmail)
	app.Get("/api/check-phone", authHandlers.HandleCheckPhone)

	authGroup := app.Group("/api", authMiddleware.JWTMiddleware)
	authGroup.Get("/users/me", authHandlers.HandleProfile)
	authGroup.Get("/users/:id", authHandlers.HandleGetUserByID(authDB.DB))
	authGroup.Get("/users", authHandlers.HandleGetAllUsers(authDB.DB))
	authGroup.Delete("/users/:id", authHandlers.HandleDeleteUser(authDB.DB))
	authGroup.Put("/users/me", authHandlers.HandleUpdateProfile)
	authGroup.Put("/users/:id/update-role", authHandlers.HandleUpdateUserRole(authDB.DB))
	app.Put("/api/users/:id/update", authHandlers.HandleUpdateProfile)
	app.Put("/api/users/:id", authHandlers.HandleUpdateProfile)

	// === MENU ===
	app.Get("/api/menu/published", menuHandlers.GetPublishedMenuItems)
	app.Get("/api/menu/published-with-category", menuHandlers.GetPublishedMenuItemsWithCategory)
	app.Get("/api/categories", menuHandlers.GetAllCategories)

	menu := app.Group("/api/menu", authMiddleware.JWTMiddleware)
	menu.Get("/", menuHandlers.GetAllMenuItems)
	menu.Get("/with-category", menuHandlers.GetAllMenuItemsWithCategory)
	menu.Get("/inventory", menuHandlers.GetInventoryItems)
	menu.Get("/calculation/:menuItemId", menuHandlers.GetCalculationByMenuItemID)
	menu.Get("/:id", menuHandlers.GetMenuItemByID)

	menu.Post("/", menuHandlers.CreateMenuItem)
	menu.Post("/inventory", menuHandlers.CreateInventoryItem)
	menu.Post("/calculation", menuHandlers.CreateCalculationForDish)
	menu.Post("/:id/publish", menuHandlers.PublishMenuItem)

	menu.Put("/:id", menuHandlers.UpdateMenuItem)
	menu.Put("/inventory/:id", menuHandlers.UpdateInventoryItem)

	menu.Delete("/:id", menuHandlers.DeleteMenuItem)
	menu.Delete("/inventory/:id", menuHandlers.DeleteInventoryItem)

	categories := app.Group("/api/categories", authMiddleware.JWTMiddleware)
	categories.Post("/", menuHandlers.CreateCategory)
	categories.Put("/:id", menuHandlers.UpdateCategory)
	categories.Delete("/:id", menuHandlers.DeleteCategory)

	// 🚀 Запуск приложения
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // По умолчанию для Koyeb
	}

	log.Printf("✅ Monolith сервер запущен на порт %s\n", port)
	log.Fatal(app.Listen(":" + port))
}



