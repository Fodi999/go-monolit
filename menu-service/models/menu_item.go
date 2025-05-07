package models

import (
	"time"

	
)

// 🥗 Меню-блюдо
type MenuItem struct {
    ID          string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description" gorm:"type:text"`
    Price       float64   `json:"price" gorm:"not null"`
    CostPrice   float64   `json:"cost_price" gorm:"not null"`
    ImageURL    string    `json:"image_url" gorm:"type:text"`           // Ссылка на картинку
    Margin      float64   `json:"margin" gorm:"not null"`               // Рассчитанная на уровне приложения
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    CategoryID  string    `json:"category_id" gorm:"type:uuid;index"`    // FK на категорию
    Published   bool      `json:"published" gorm:"default:false;index"`  // Опубликовано ли блюдо
}

// 📂 Меню-блюдо с категорией (JOIN)
type MenuItemWithCategory struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	CostPrice    float64   `json:"cost_price"`
	ImageURL     string    `json:"image_url"`
	Margin       float64   `json:"margin"`
	CreatedAt    time.Time `json:"created_at"`
	CategoryID   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Published    bool      `json:"published"`
}

// 📦 Продукт на складе
type InventoryItem struct {
	ID          string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductName string    `json:"product_name" gorm:"not null"`
	WeightGrams int       `json:"weight_grams" gorm:"not null"`
	PricePerKg  float64   `json:"price_per_kg" gorm:"not null"`
	Available   bool      `json:"available" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	Emoji       string    `json:"emoji" gorm:"default:'🍽️'"`
	Category    *string   `json:"category" gorm:"default:'прочее'"`
}

// 📐 Ингредиент в калькуляции
type CalculationIngredient struct {
	ID              string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CalculationID   string    `json:"calculation_id" gorm:"type:uuid;not null;index"`
	ProductName     string    `json:"product_name" gorm:"not null"`
	AmountGrams     int       `json:"amount_grams" gorm:"not null"`
	PricePerKg      float64   `json:"price_per_kg" gorm:"not null"`
	WastePercent    float64   `json:"waste_percent" gorm:"default:0.0"`
	PriceAfterWaste float64   `json:"price_after_waste"`
	TotalCost       float64   `json:"total_cost"`
	CreatedAt       time.Time `json:"created_at"`
}

// 🧾 Финальная калькуляция блюда
type MenuCalculation struct {
	ID               string                  `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MenuItemID       string                  `json:"menu_item_id" gorm:"type:uuid;not null;index"`
	TotalOutputGrams int                     `json:"total_output_grams"`
	TotalCost        float64                 `json:"total_cost"`
	Ingredients      []CalculationIngredient `json:"ingredients" gorm:"foreignKey:CalculationID;constraint:OnDelete:CASCADE"`
	CreatedAt        time.Time               `json:"created_at"`
}

// ✅ Псевдоним для совместимости
type Calculation = MenuCalculation



