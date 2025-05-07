package utils

import (
	"monolith/menu-service/models"
	"monolith/menu-service/utils/emoji"

	"gorm.io/gorm"
)

// Получить все продукты со склада
func GetAllInventoryItems(db *gorm.DB) ([]models.InventoryItem, error) {
	var items []models.InventoryItem
	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Создать продукт на складе
func CreateInventoryItem(db *gorm.DB, item *models.InventoryItem) error {
	item.Emoji = emoji.GenerateEmoji(item.ProductName)
	category := emoji.GenerateCategory(item.ProductName)
	item.Category = &category

	return db.Create(item).Error
}

// Обновить продукт на складе
func UpdateInventoryItem(db *gorm.DB, id string, item *models.InventoryItem) error {
	item.Emoji = emoji.GenerateEmoji(item.ProductName)
	category := emoji.GenerateCategory(item.ProductName)
	item.Category = &category

	var existing models.InventoryItem
	if err := db.First(&existing, "id = ?", id).Error; err != nil {
		return err
	}

	existing.ProductName = item.ProductName
	existing.WeightGrams = item.WeightGrams
	existing.PricePerKg = item.PricePerKg
	existing.Available = item.Available
	existing.Emoji = item.Emoji
	existing.Category = item.Category

	return db.Save(&existing).Error
}

// Удалить продукт со склада
func DeleteInventoryItem(db *gorm.DB, id string) error {
	return db.Delete(&models.InventoryItem{}, "id = ?", id).Error
}




