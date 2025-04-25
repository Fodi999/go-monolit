package utils

import (
	"database/sql"
	"monolith/menu-service/models"
	"monolith/menu-service/utils/emoji"
	"time"
)

// Получить все продукты со склада
func GetAllInventoryItems(db *sql.DB) ([]models.InventoryItem, error) {
	rows, err := db.Query(`
		SELECT id, product_name, weight_grams, price_per_kg, available, created_at, emoji
		FROM inventory_items
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		err := rows.Scan(
			&item.ID,
			&item.ProductName,
			&item.WeightGrams,
			&item.PricePerKg,
			&item.Available,
			&item.CreatedAt,
			&item.Emoji,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// Создать продукт на складе
func CreateInventoryItem(db *sql.DB, item *models.InventoryItem) error {
	// ✅ Всегда генерируем emoji автоматически на основе названия
	item.Emoji = emoji.GenerateEmoji(item.ProductName)

	query := `
		INSERT INTO inventory_items (product_name, weight_grams, price_per_kg, available, created_at, emoji)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return db.QueryRow(query,
		item.ProductName,
		item.WeightGrams,
		item.PricePerKg,
		item.Available,
		time.Now(),
		item.Emoji,
	).Scan(&item.ID)
}

// Обновить продукт на складе
func UpdateInventoryItem(db *sql.DB, id string, item *models.InventoryItem) error {
	// ✅ При обновлении тоже можно пересчитать emoji
	item.Emoji = emoji.GenerateEmoji(item.ProductName)

	_, err := db.Exec(`
		UPDATE inventory_items 
		SET product_name=$1, weight_grams=$2, price_per_kg=$3, available=$4, emoji=$5
		WHERE id=$6
	`,
		item.ProductName,
		item.WeightGrams,
		item.PricePerKg,
		item.Available,
		item.Emoji,
		id,
	)
	return err
}

// Удалить продукт со склада
func DeleteInventoryItem(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM inventory_items WHERE id=$1`, id)
	return err
}


