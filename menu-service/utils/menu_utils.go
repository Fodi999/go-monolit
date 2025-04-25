// C:\Users\Admin\Desktop\fo-sushi\backend\menu-service\utils\menu_utils.go
package utils

import (
	"database/sql"
	"monolith/menu-service/models"
	"time"
)

// 🥗 Получить все блюда без категории
func GetAllMenuItems(db *sql.DB) ([]models.MenuItem, error) {
	rows, err := db.Query(`
		SELECT id, name, description, price, cost_price, image_url, margin, created_at, category_id, published
		FROM menu_items
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.CostPrice,
			&item.ImageURL,
			&item.Margin,
			&item.CreatedAt,
			&item.CategoryID,
			&item.Published,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// 📂 Получить все блюда с категорией (JOIN)
func GetAllMenuItemsWithCategory(db *sql.DB) ([]models.MenuItemWithCategory, error) {
	const q = `
		SELECT
			m.id, m.name, m.description, m.price, m.cost_price,
			m.image_url, m.margin, m.created_at, m.category_id,
			COALESCE(c.name, '') AS category_name,
			m.published
		FROM menu_items m
		LEFT JOIN categories c ON m.category_id = c.id
		ORDER BY m.created_at DESC
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItemWithCategory
	for rows.Next() {
		var item models.MenuItemWithCategory
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.CostPrice,
			&item.ImageURL,
			&item.Margin,
			&item.CreatedAt,
			&item.CategoryID,
			&item.CategoryName,
			&item.Published,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// 🔍 Получить блюдо по ID
func GetMenuItemByID(db *sql.DB, id string) (*models.MenuItem, error) {
	const q = `
		SELECT id, name, description, price, cost_price, image_url, margin, created_at, category_id, published
		FROM menu_items
		WHERE id = $1
	`
	row := db.QueryRow(q, id)

	var item models.MenuItem
	if err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.CostPrice,
		&item.ImageURL,
		&item.Margin,
		&item.CreatedAt,
		&item.CategoryID,
		&item.Published,
	); err != nil {
		return nil, err
	}
	return &item, nil
}

// ➕ Создание блюда
func CreateMenuItem(db *sql.DB, item *models.MenuItem) error {
	item.Margin = item.Price - item.CostPrice
	const q = `
		INSERT INTO menu_items
			(name, description, price, cost_price, image_url, margin, created_at, category_id, published)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`
	return db.QueryRow(q,
		item.Name,
		item.Description,
		item.Price,
		item.CostPrice,
		item.ImageURL,
		item.Margin,
		time.Now(),
		item.CategoryID,
		item.Published, // по умолчанию false
	).Scan(&item.ID)
}

// ✏️ Обновление блюда (не трогаем published)
func UpdateMenuItem(db *sql.DB, id string, item *models.MenuItem) error {
	item.Margin = item.Price - item.CostPrice
	const q = `
		UPDATE menu_items
		SET name=$1, description=$2, price=$3, cost_price=$4,
		    image_url=$5, margin=$6, category_id=$7
		WHERE id=$8
	`
	_, err := db.Exec(q,
		item.Name,
		item.Description,
		item.Price,
		item.CostPrice,
		item.ImageURL,
		item.Margin,
		item.CategoryID,
		id,
	)
	return err
}

// 🔄 Переключить статус публикации блюда
func TogglePublishMenuItem(db *sql.DB, id string) error {
	const q = `
		UPDATE menu_items
		SET published = NOT published
		WHERE id = $1
	`
	_, err := db.Exec(q, id)
	return err
}

// 📥 Получить все опубликованные блюда
func GetPublishedMenuItems(db *sql.DB) ([]models.MenuItem, error) {
	const q = `
		SELECT id, name, description, price, cost_price, image_url, margin, created_at, category_id, published
		FROM menu_items
		WHERE published = TRUE
		ORDER BY created_at DESC
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.CostPrice,
			&item.ImageURL,
			&item.Margin,
			&item.CreatedAt,
			&item.CategoryID,
			&item.Published,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// 🧮 Сохранение или обновление калькуляции блюда
func SaveDishCalculation(db *sql.DB, calc *models.Calculation) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, есть ли запись
	var existingID string
	err = tx.QueryRow(`
		SELECT id
		FROM menu_calculations
		WHERE menu_item_id = $1
		LIMIT 1
	`, calc.MenuItemID).Scan(&existingID)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if existingID != "" {
		// обновляем
		_, err = tx.Exec(`
			UPDATE menu_calculations
			SET total_output_grams=$1, total_cost=$2, created_at=$3
			WHERE id=$4
		`, calc.TotalOutputGrams, calc.TotalCost, time.Now(), existingID)
		if err != nil {
			return err
		}
		// чистим старые ингредиенты
		_, err = tx.Exec(`DELETE FROM calculation_ingredients WHERE calculation_id = $1`, existingID)
		if err != nil {
			return err
		}
		calc.ID = existingID
	} else {
		// создаём новую
		err = tx.QueryRow(`
			INSERT INTO menu_calculations
				(menu_item_id, total_output_grams, total_cost, created_at)
			VALUES ($1,$2,$3,$4)
			RETURNING id
		`, calc.MenuItemID, calc.TotalOutputGrams, calc.TotalCost, time.Now()).Scan(&calc.ID)
		if err != nil {
			return err
		}
	}

	// вставляем ингредиенты
	for _, ing := range calc.Ingredients {
		_, err := tx.Exec(`
			INSERT INTO calculation_ingredients
				(calculation_id, product_name, amount_grams, price_per_kg,
				 waste_percent, price_after_waste, total_cost, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`, calc.ID,
			ing.ProductName,
			ing.AmountGrams,
			ing.PricePerKg,
			ing.WastePercent,
			ing.PriceAfterWaste,
			ing.TotalCost,
			time.Now(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// 📊 Получить калькуляцию по блюду
func GetCalculationByMenuItemID(db *sql.DB, menuItemID string) (*models.Calculation, error) {
	const qCalc = `
		SELECT id, menu_item_id, total_output_grams, total_cost, created_at
		FROM menu_calculations
		WHERE menu_item_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var calc models.Calculation
	if err := db.QueryRow(qCalc, menuItemID).Scan(
		&calc.ID,
		&calc.MenuItemID,
		&calc.TotalOutputGrams,
		&calc.TotalCost,
		&calc.CreatedAt,
	); err != nil {
		return nil, err
	}

	// ингредиенты
	rows, err := db.Query(`
		SELECT product_name, amount_grams, price_per_kg,
		       waste_percent, price_after_waste, total_cost
		FROM calculation_ingredients
		WHERE calculation_id = $1
	`, calc.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ing models.CalculationIngredient
		if err := rows.Scan(
			&ing.ProductName,
			&ing.AmountGrams,
			&ing.PricePerKg,
			&ing.WastePercent,
			&ing.PriceAfterWaste,
			&ing.TotalCost,
		); err != nil {
			return nil, err
		}
		calc.Ingredients = append(calc.Ingredients, ing)
	}

	return &calc, nil
}

// ❌ Удаление блюда
func DeleteMenuItem(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM menu_items WHERE id = $1`, id)
	return err
}







