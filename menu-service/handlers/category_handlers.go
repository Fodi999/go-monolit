package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"monolith/menu-service/database"
	"monolith/menu-service/models"
)

// Получить все категории
func GetAllCategories(c *fiber.Ctx) error {
	rows, err := database.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при получении категорий",
		})
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при сканировании категории",
			})
		}
		categories = append(categories, cat)
	}

	return c.JSON(categories)
}

// Создать новую категорию
func CreateCategory(c *fiber.Ctx) error {
	var input models.Category
	if err := json.Unmarshal(c.Body(), &input); err != nil || input.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Невалидный ввод. Поле name обязательно",
		})
	}

	query := `INSERT INTO categories (name) VALUES ($1) RETURNING id`
	err := database.DB.QueryRow(query, input.Name).Scan(&input.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка при создании категории",
		})
	}

	return c.Status(http.StatusCreated).JSON(input)
}

// Обновить категорию
func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var input models.Category

	if err := c.BodyParser(&input); err != nil || input.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Неверный ввод. Поле name обязательно",
		})
	}

	// Проверка существования
	var exists bool
	err := database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1)`, id).Scan(&exists)
	if err != nil || !exists {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Категория не найдена",
		})
	}

	_, err = database.DB.Exec(`UPDATE categories SET name=$1 WHERE id=$2`, input.Name, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось обновить категорию",
		})
	}

	return c.JSON(fiber.Map{"message": "Категория обновлена"})
}

// Удалить категорию
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	// Проверка существования
	var exists bool
	err := database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1)`, id).Scan(&exists)
	if err != nil || !exists {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Категория не найдена",
		})
	}

	_, err = database.DB.Exec(`DELETE FROM categories WHERE id=$1`, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Не удалось удалить категорию",
		})
	}

	return c.JSON(fiber.Map{"message": "Категория удалена"})
}


