package utils

import (
	"monolith/menu-service/models"
	"time"

	"gorm.io/gorm"
)

// Получить все блюда без категории
func GetAllMenuItems(db *gorm.DB) ([]models.MenuItem, error) {
	var items []models.MenuItem
	if err := db.Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Получить все блюда с категорией
func GetAllMenuItemsWithCategory(db *gorm.DB) ([]models.MenuItemWithCategory, error) {
	var items []models.MenuItemWithCategory
	err := db.Table("menu_items").
		Select("menu_items.*, categories.name as category_name").
		Joins("LEFT JOIN categories ON menu_items.category_id = categories.id").
		Order("menu_items.created_at DESC").
		Scan(&items).Error
	return items, err
}

// Получить блюдо по ID
func GetMenuItemByID(db *gorm.DB, id string) (*models.MenuItem, error) {
	var item models.MenuItem
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Создать блюдо
func CreateMenuItem(db *gorm.DB, item *models.MenuItem) error {
	item.Margin = item.Price - item.CostPrice
	item.CreatedAt = time.Now()
	return db.Create(item).Error
}

// Обновить блюдо
func UpdateMenuItem(db *gorm.DB, id string, input *models.MenuItem) error {
	var item models.MenuItem
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		return err
	}
	input.Margin = input.Price - input.CostPrice
	return db.Model(&item).Updates(input).Error
}

// Переключить публикацию
func TogglePublishMenuItem(db *gorm.DB, id string) error {
	var item models.MenuItem
	if err := db.First(&item, "id = ?", id).Error; err != nil {
		return err
	}
	item.Published = !item.Published
	return db.Save(&item).Error
}

// Получить опубликованные блюда
func GetPublishedMenuItems(db *gorm.DB) ([]models.MenuItem, error) {
	var items []models.MenuItem
	if err := db.Where("published = TRUE").Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Удалить блюдо
func DeleteMenuItem(db *gorm.DB, id string) error {
	return db.Delete(&models.MenuItem{}, "id = ?", id).Error
}

// Сохранить или обновить калькуляцию
func SaveDishCalculation(db *gorm.DB, calc *models.Calculation) error {
	var existing models.Calculation
	if err := db.Where("menu_item_id = ?", calc.MenuItemID).First(&existing).Error; err == nil {
		calc.ID = existing.ID
		calc.CreatedAt = time.Now()
		if err := db.Model(&existing).Updates(calc).Error; err != nil {
			return err
		}
		_ = db.Where("calculation_id = ?", existing.ID).Delete(&models.CalculationIngredient{})
	} else {
		calc.CreatedAt = time.Now()
		if err := db.Create(calc).Error; err != nil {
			return err
		}
	}

	for _, ing := range calc.Ingredients {
		ing.CalculationID = calc.ID
		ing.CreatedAt = time.Now()
		if err := db.Create(&ing).Error; err != nil {
			return err
		}
	}
	return nil
}

// Получить калькуляцию по блюду
func GetCalculationByMenuItemID(db *gorm.DB, menuItemID string) (*models.Calculation, error) {
	var calc models.Calculation
	if err := db.
		Where("menu_item_id = ?", menuItemID).
		Order("created_at DESC").
		First(&calc).Error; err != nil {
		return nil, err
	}

	var ingredients []models.CalculationIngredient
	if err := db.Where("calculation_id = ?", calc.ID).
		Find(&ingredients).Error; err != nil {
		return nil, err
	}
	calc.Ingredients = ingredients

	return &calc, nil
}








