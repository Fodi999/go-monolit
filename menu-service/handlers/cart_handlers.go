package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"monolith/menu-service/database"
	"monolith/menu-service/models"
)

// Тело запроса на добавление товара
type AddToCartBody struct {
	MenuItemID string  `json:"menuItemId"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
}

// 📦 Получить корзину пользователя
func GetCart(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var cart models.Cart
	if err := database.DB.
		Preload("Items").
		FirstOrCreate(&cart, models.Cart{UserID: userID}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "не удалось получить корзину"})
	}

	// Подтягиваем ImageURL для каждой позиции
	for i := range cart.Items {
		var mi models.MenuItem
		if err := database.DB.
			Select("image_url").
			First(&mi, "id = ?", cart.Items[i].MenuItemID).
			Error; err == nil {
			cart.Items[i].ImageURL = mi.ImageURL
		}
	}

	return c.JSON(cart)
}

// ➕ Добавить товар в корзину
func AddToCart(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var body AddToCartBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "неверный JSON"})
	}

	var cart models.Cart
	if err := database.DB.FirstOrCreate(&cart, models.Cart{UserID: userID}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "не удалось найти или создать корзину"})
	}

	var item models.CartItem
	err := database.DB.
		Where("cart_id = ? AND menu_item_id = ?", cart.ID, body.MenuItemID).
		First(&item).Error

	if err == gorm.ErrRecordNotFound {
		item = models.CartItem{
			CartID:     cart.ID,
			MenuItemID: body.MenuItemID,
			Name:       body.Name,
			Quantity:   body.Quantity,
			Price:      body.Price,
		}
		if err := database.DB.Create(&item).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "не удалось добавить товар"})
		}
	} else if err == nil {
		item.Quantity += body.Quantity
		if err := database.DB.Save(&item).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "не удалось обновить товар"})
		}
	} else {
		return c.Status(500).JSON(fiber.Map{"error": "ошибка обработки корзины"})
	}

	// Подтягиваем картинку перед ответом
	var mi models.MenuItem
	if err := database.DB.Select("image_url").
		First(&mi, "id = ?", item.MenuItemID).Error; err == nil {
		item.ImageURL = mi.ImageURL
	}

	return c.JSON(item)
}

// ✏️ Обновить количество товара
func UpdateCartItem(c *fiber.Ctx) error {
	userID := c.Params("userId")
	menuItemID := c.Params("menuItemId")

	var body struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "неверный JSON"})
	}

	var cart models.Cart
	if err := database.DB.First(&cart, "user_id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "корзина не найдена"})
	}

	var item models.CartItem
	if err := database.DB.
		Where("cart_id = ? AND menu_item_id = ?", cart.ID, menuItemID).
		First(&item).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "товар не найден"})
	}

	if body.Quantity > 0 {
		item.Quantity = body.Quantity
		if err := database.DB.Save(&item).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "не удалось обновить товар"})
		}
	} else {
		if err := database.DB.Delete(&item).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "не удалось удалить товар"})
		}
		return c.SendStatus(204)
	}

	// Подтягиваем картинку перед ответом
	var mi models.MenuItem
	if err := database.DB.Select("image_url").
		First(&mi, "id = ?", item.MenuItemID).Error; err == nil {
		item.ImageURL = mi.ImageURL
	}

	return c.JSON(item)
}

// 🚮 Удалить одну позицию из корзины
func RemoveCartItem(c *fiber.Ctx) error {
	userID := c.Params("userId")
	menuItemID := c.Params("menuItemId")

	var cart models.Cart
	if err := database.DB.First(&cart, "user_id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "корзина не найдена"})
	}

	if err := database.DB.
		Where("cart_id = ? AND menu_item_id = ?", cart.ID, menuItemID).
		Delete(&models.CartItem{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "не удалось удалить товар"})
	}

	return c.SendStatus(204)
}

// 🧹 Очистить всю корзину
func ClearCart(c *fiber.Ctx) error {
	userID := c.Params("userId")

	var cart models.Cart
	if err := database.DB.First(&cart, "user_id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "корзина не найдена"})
	}

	if err := database.DB.
		Where("cart_id = ?", cart.ID).
		Delete(&models.CartItem{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "не удалось очистить корзину"})
	}

	return c.SendStatus(204)
}




