package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	// Auth
	authDB "monolith/auth_service/database"
	authHandlers "monolith/auth_service/handlers"
	authMiddleware "monolith/auth_service/middleware"

	// Menu (включая корзину и заказы)
	menuDB "monolith/menu-service/database"
	menuMiddleware "monolith/menu-service/middleware"
	menuRoutes "monolith/menu-service/routes"
)

func main() {
	// === Загрузка переменных окружения ===
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env файл не найден — используем переменные окружения")
	}

	// === Подключение к базе MENU (через GORM) ===
	menuDSN := os.Getenv("MENU_DATABASE_URL")
	db, err := gorm.Open(postgres.Open(menuDSN), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Warn), // Уровень: Silent, Error, Warn, Info
	})
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к базе MENU: %v", err)
	}
	menuDB.Init(db)
	log.Println("✅ Подключение и миграция базы MENU успешно")

	// === Подключение к базе AUTH ===
	authDB.Init()
	log.Println("✅ Подключение к базе AUTH успешно")

	// === Инициализация Fiber ===
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// === МИДЛВАРЫ ===
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(fiberLogger.New())
	app.Use(menuMiddleware.LoggerMiddleware())

	// === Базовые эндпоинты ===
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("✅ Backend is up and running")
	})
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	// === AUTH ===
	app.Post("/api/register", authHandlers.HandleRegister)
	app.Post("/api/login", authHandlers.HandleLogin)
	app.Get("/api/check-email", authHandlers.HandleCheckEmail)
	app.Get("/api/check-phone", authHandlers.HandleCheckPhone)

	// === Защищённые /api маршруты ===
	api := app.Group("/api", authMiddleware.JWTMiddleware)
	api.Get("/users/me", authHandlers.HandleProfile)
	api.Get("/users/:id", authHandlers.HandleGetUserByID(authDB.DB))
	api.Get("/users", authHandlers.HandleGetAllUsers(authDB.DB))
	api.Delete("/users/:id", authHandlers.HandleDeleteUser(authDB.DB))
	api.Put("/users/me", authHandlers.HandleUpdateProfile)
	api.Put("/users/:id/update-role", authHandlers.HandleUpdateUserRole(authDB.DB))

	// === Роуты меню и категории ===
	menuRoutes.SetupMenuRoutes(app)
	menuRoutes.SetupCategoryRoutes(app)

	// === Роуты корзины и заказов ===
	menuRoutes.SetupCartRoutes(app)
	menuRoutes.SetupOrderRoutes(app)

	// === Запуск сервера ===
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Printf("🚀 Monolith сервер запущен на порту %s\n", port)
	log.Fatal(app.Listen(":" + port))
}





