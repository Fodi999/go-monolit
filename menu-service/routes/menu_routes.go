package routes

import (
    "github.com/gofiber/fiber/v2"
    "monolith/menu-service/handlers"
)

// Настройка основных маршрутов меню и склада
func SetupMenuRoutes(app *fiber.App) {
    menu := app.Group("/api/menu")

    // Склад
    menu.Get("/inventory", handlers.GetInventoryItems)
    menu.Post("/inventory", handlers.CreateInventoryItem)
    menu.Put("/inventory/:id", handlers.UpdateInventoryItem)
    menu.Delete("/inventory/:id", handlers.DeleteInventoryItem)

    // Калькуляция блюда
    menu.Get("/calculation/:menuItemId", handlers.GetCalculationByMenuItemID)
    menu.Post("/calculation", handlers.CreateCalculationForDish)

    // Публичное меню
    menu.Get("/published", handlers.GetPublishedMenuItems)
    menu.Get("/published-with-category", handlers.GetPublishedMenuItemsWithCategory)

    // Администрирование меню
    menu.Get("/with-category", handlers.GetAllMenuItemsWithCategory)
    menu.Get("/", handlers.GetAllMenuItems)
    menu.Get("/:id", handlers.GetMenuItemByID)
    menu.Post("/", handlers.CreateMenuItem)
    menu.Put("/:id", handlers.UpdateMenuItem)
    menu.Delete("/:id", handlers.DeleteMenuItem)

    // Публикация / снятие публикации
    menu.Post("/:id/publish", handlers.PublishMenuItem)
}

// Маршруты корзины
// SetupCartRoutes регистрирует маршруты корзины
func SetupCartRoutes(app *fiber.App) {
    cart := app.Group("/api/users/:userId/cart")
    cart.Get("/",    handlers.GetCart)
    cart.Post("/",   handlers.AddToCart)
    cart.Put("/:menuItemId", handlers.UpdateCartItem)
    cart.Delete("/:menuItemId", handlers.RemoveCartItem)
    cart.Delete("/", handlers.ClearCart)
}



// Маршруты заказов
// menu-service/routes/menu_routes.go
func SetupOrderRoutes(app *fiber.App) {
    order := app.Group("/api/users/:userId")
    order.Post("/order", handlers.PlaceOrder)      // ← PlaceOrder, а не PlaceOrderHandler
    order.Get("/orders", handlers.GetUserOrders)   // ← GetUserOrders, а не GetUserOrdersHandler
}








