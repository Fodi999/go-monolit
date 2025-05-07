package handlers

import (
	"github.com/gofiber/fiber/v2"
	"monolith/menu-service/database"
	"monolith/menu-service/models"
)

func PlaceOrder(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var cart models.Cart
	if err := database.DB.Preload("Items").First(&cart, "user_id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "корзина не найдена"})
	}

	if len(cart.Items) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "корзина пуста"})
	}

	var total float64
	var orderItems []models.OrderItem
	for _, item := range cart.Items {
		total += float64(item.Quantity) * item.Price
		orderItems = append(orderItems, models.OrderItem{
			MenuItemID: item.MenuItemID,
			Name:       item.Name,
			Quantity:   item.Quantity,
			Price:      item.Price,
		})
	}

	order := models.Order{
		UserID:     userID,
		CartID:     cart.ID,
		Items:      orderItems,
		TotalPrice: total,
		Status:     "pending",
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "не удалось создать заказ"})
	}

	// Очистить корзину
	database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})

	return c.Status(201).JSON(order)
}

func GetUserOrders(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var orders []models.Order
	if err := database.DB.
		Preload("Items").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "не удалось получить заказы"})
	}

	return c.JSON(orders)
}
