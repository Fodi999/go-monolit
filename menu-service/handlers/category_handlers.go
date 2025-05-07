package handlers

import (
	"github.com/gofiber/fiber/v2"
	"monolith/menu-service/database"
	"monolith/menu-service/models"
)

// 📥 Получить все категории
func GetAllCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при получении категорий",
		})
	}
	return c.JSON(categories)
}

// ➕ Создать категорию
func CreateCategory(c *fiber.Ctx) error {
	var input models.Category
	if err := c.BodyParser(&input); err != nil || input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Невалидный ввод. Поле name обязательно",
		})
	}

	if err := database.DB.Create(&input).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при создании категории",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(input)
}

// ✏️ Обновить категорию
func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var input models.Category
	if err := c.BodyParser(&input); err != nil || input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ввод. Поле name обязательно",
		})
	}

	var category models.Category
	if err := database.DB.First(&category, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Категория не найдена",
		})
	}

	category.Name = input.Name
	if err := database.DB.Save(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось обновить категорию",
		})
	}

	return c.JSON(fiber.Map{"message": "Категория обновлена"})
}

// ❌ Удалить категорию
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var category models.Category
	if err := database.DB.First(&category, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Категория не найдена",
		})
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось удалить категорию",
		})
	}

	return c.JSON(fiber.Map{"message": "Категория удалена"})
}



