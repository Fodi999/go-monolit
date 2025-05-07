package database

import (
	"gorm.io/gorm"
	"monolith/menu-service/models"
)

var DB *gorm.DB

func Init(db *gorm.DB) {
	DB = db

	DB.AutoMigrate(
		&models.Category{},
		&models.MenuItem{},
		&models.InventoryItem{},
		&models.Calculation{},
		&models.CalculationIngredient{}, // 🔥 Добавляем!
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	)
}


