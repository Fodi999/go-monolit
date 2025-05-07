package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"monolith/menu-service/database"
	"monolith/menu-service/models"
)

// 📦 Получить все блюда без категории
func GetAllMenuItems(c *fiber.Ctx) error {
	var items []models.MenuItem
	if err := database.DB.Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить блюда",
		})
	}
	return c.JSON(items)
}

// 📂 Получить все опубликованные блюда с названием категории
func GetPublishedMenuItemsWithCategory(c *fiber.Ctx) error {
	var result []models.MenuItemWithCategory

	err := database.DB.Table("menu_items").
		Select("menu_items.*, categories.name as category_name").
		Joins("LEFT JOIN categories ON menu_items.category_id = categories.id").
		Where("menu_items.published = TRUE").
		Scan(&result).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить опубликованные блюда с категориями",
		})
	}

	return c.JSON(result)
}

// 📂 Получить все блюда с названием категории (JOIN)
func GetAllMenuItemsWithCategory(c *fiber.Ctx) error {
	var result []models.MenuItemWithCategory

	err := database.DB.Table("menu_items").
		Select("menu_items.*, categories.name as category_name").
		Joins("LEFT JOIN categories ON menu_items.category_id = categories.id").
		Scan(&result).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить блюда с категориями",
		})
	}

	return c.JSON(result)
}

// 📦 Получить все опубликованные блюда
func GetPublishedMenuItems(c *fiber.Ctx) error {
	var items []models.MenuItem
	if err := database.DB.Where("published = TRUE").Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить опубликованные блюда",
		})
	}
	return c.JSON(items)
}

// 📊 Получить калькуляцию по блюду
func GetCalculationByMenuItemID(c *fiber.Ctx) error {
	menuItemID := c.Params("menuItemId")
	var calc models.Calculation

	if err := database.DB.
		Where("menu_item_id = ?", menuItemID).
		First(&calc).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Калькуляция не найдена",
		})
	}
	return c.JSON(calc)
}

// 🍽 Получить блюдо по ID
func GetMenuItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.MenuItem

	if err := database.DB.First(&item, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Блюдо не найдено",
		})
	}
	return c.JSON(item)
}

// 🍽 Создать новое блюдо
func CreateMenuItem(c *fiber.Ctx) error {
	var item models.MenuItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат запроса",
		})
	}
	item.Margin = item.Price - item.CostPrice

	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось сохранить блюдо",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

// 🛠 Обновить блюдо
func UpdateMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var input models.MenuItem

	if err := json.Unmarshal(c.Body(), &input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат тела запроса",
		})
	}
	input.Margin = input.Price - input.CostPrice

	var item models.MenuItem
	if err := database.DB.First(&item, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Блюдо не найдено",
		})
	}

	item.Name = input.Name
	item.Description = input.Description
	item.Price = input.Price
	item.CostPrice = input.CostPrice
	item.Margin = input.Margin
	item.ImageURL = input.ImageURL
	item.CategoryID = input.CategoryID

	if err := database.DB.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось обновить блюдо",
		})
	}

	return c.JSON(item)
}

// 🔄 Опубликовать/снять с публикации блюдо
func PublishMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.MenuItem

	if err := database.DB.First(&item, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Блюдо не найдено",
		})
	}

	item.Published = !item.Published
	if err := database.DB.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось изменить статус публикации блюда",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// ❌ Удалить блюдо
func DeleteMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.MenuItem{}, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось удалить блюдо",
		})
	}
	return c.JSON(fiber.Map{"message": "Блюдо удалено"})
}

// 📦 Получить все продукты со склада
func GetInventoryItems(c *fiber.Ctx) error {
	var items []models.InventoryItem
	if err := database.DB.Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить складские продукты",
		})
	}
	return c.JSON(items)
}

// ➕ Добавить продукт на склад
func CreateInventoryItem(c *fiber.Ctx) error {
	var item models.InventoryItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат тела запроса",
		})
	}
	if err := database.DB.Create(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось добавить продукт на склад",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

// 🛠 Обновить продукт на складе
func UpdateInventoryItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var input models.InventoryItem

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат тела запроса",
		})
	}

	var item models.InventoryItem
	if err := database.DB.First(&item, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Продукт не найден",
		})
	}

	item.ProductName = input.ProductName
	item.WeightGrams = input.WeightGrams
	item.PricePerKg = input.PricePerKg
	item.Available = input.Available
	item.Emoji = input.Emoji
	item.Category = input.Category

	if err := database.DB.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось обновить продукт",
		})
	}
	return c.JSON(item)
}

// 🗑 Удалить продукт со склада
func DeleteInventoryItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := database.DB.Delete(&models.InventoryItem{}, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось удалить продукт со склада",
		})
	}
	return c.JSON(fiber.Map{"message": "Продукт удалён"})
}

// 📊 Сохранить калькуляцию
func CreateCalculationForDish(c *fiber.Ctx) error {
	var calc models.Calculation
	if err := c.BodyParser(&calc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат тела запроса",
		})
	}
	if err := database.DB.Create(&calc).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось сохранить калькуляцию",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Калькуляция сохранена"})
}








