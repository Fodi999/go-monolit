package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"monolith/menu-service/database"
	"monolith/menu-service/models"
	"monolith/menu-service/utils"
)

// 📦 Получить все блюда без категории
func GetAllMenuItems(c *fiber.Ctx) error {
	items, err := utils.GetAllMenuItems(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить блюда",
		})
	}
	if items == nil {
		items = []models.MenuItem{}
	}
	return c.JSON(items)
}
// 📂 Получить все опубликованные блюда с названием категории
func GetPublishedMenuItemsWithCategory(c *fiber.Ctx) error {
	items, err := utils.GetAllMenuItemsWithCategory(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить опубликованные блюда с категориями",
		})
	}

	// Фильтруем только опубликованные
	var published []models.MenuItemWithCategory
	for _, item := range items {
		if item.Published {
			published = append(published, item)
		}
	}

	return c.JSON(published)
}

// 📂 Получить все блюда с названием категории (JOIN)
func GetAllMenuItemsWithCategory(c *fiber.Ctx) error {
	items, err := utils.GetAllMenuItemsWithCategory(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить блюда с категориями",
		})
	}
	if items == nil {
		items = []models.MenuItemWithCategory{}
	}
	return c.JSON(items)
}

// 📦 Получить все опубликованные блюда
func GetPublishedMenuItems(c *fiber.Ctx) error {
	items, err := utils.GetPublishedMenuItems(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить опубликованные блюда",
		})
	}
	if items == nil {
		items = []models.MenuItem{}
	}
	return c.JSON(items)
}

// 📊 Получить калькуляцию по блюду
func GetCalculationByMenuItemID(c *fiber.Ctx) error {
	menuItemID := c.Params("menuItemId")
	calc, err := utils.GetCalculationByMenuItemID(database.DB, menuItemID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Калькуляция не найдена",
		})
	}
	return c.JSON(calc)
}

// 🍽 Получить блюдо по ID
func GetMenuItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	item, err := utils.GetMenuItemByID(database.DB, id)
	if err != nil {
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
	if err := utils.CreateMenuItem(database.DB, &item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось сохранить блюдо",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

// 🛠 Обновить блюдо
func UpdateMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.MenuItem
	if err := json.Unmarshal(c.Body(), &item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат тела запроса",
		})
	}
	item.Margin = item.Price - item.CostPrice
	if err := utils.UpdateMenuItem(database.DB, id, &item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось обновить блюдо",
		})
	}
	return c.JSON(item)
}

// 🔄 Опубликовать (или снять публикацию) блюда
func PublishMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.TogglePublishMenuItem(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось изменить статус публикации блюда",
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

// ❌ Удалить блюдо
func DeleteMenuItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.DeleteMenuItem(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось удалить блюдо",
		})
	}
	return c.JSON(fiber.Map{"message": "Блюдо удалено"})
}

// 📦 Получить все продукты со склада
func GetInventoryItems(c *fiber.Ctx) error {
	items, err := utils.GetAllInventoryItems(database.DB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось получить складские продукты",
		})
	}
	if items == nil {
		items = []models.InventoryItem{}
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
	if err := utils.CreateInventoryItem(database.DB, &item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось добавить продукт на склад",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(item)
}

// 🛠 Обновить продукт на складе
func UpdateInventoryItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.InventoryItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат тела запроса",
		})
	}
	if err := utils.UpdateInventoryItem(database.DB, id, &item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось обновить продукт",
		})
	}
	return c.JSON(item)
}

// 🗑 Удалить продукт со склада
func DeleteInventoryItem(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := utils.DeleteInventoryItem(database.DB, id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось удалить продукт со склада",
		})
	}
	return c.JSON(fiber.Map{"message": "Продукт удалён"})
}

// 📊 Сохранить калькуляционную карту блюда
func CreateCalculationForDish(c *fiber.Ctx) error {
	var calc models.Calculation
	if err := c.BodyParser(&calc); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный формат тела запроса",
		})
	}
	if err := utils.SaveDishCalculation(database.DB, &calc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось сохранить калькуляцию",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Калькуляция сохранена"})
}







