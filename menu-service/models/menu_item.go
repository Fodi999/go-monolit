package models

import "time"

// 🥗 Меню-блюдо
type MenuItem struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CostPrice   float64   `json:"cost_price"`
	ImageURL    string    `json:"image_url"`
	Margin      float64   `json:"margin"`
	CreatedAt   time.Time `json:"created_at"`
	CategoryID  string    `json:"category_id"`
	Published   bool      `json:"published"` // новое поле
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
	Published    bool      `json:"published"` // и здесь
}

// 📦 Продукт на складе
type InventoryItem struct {
	ID          string    `json:"id"`
	ProductName string    `json:"product_name"`
	WeightGrams int       `json:"weight_grams"`
	PricePerKg  float64   `json:"price_per_kg"`
	Available   bool      `json:"available"`
	CreatedAt   time.Time `json:"created_at"`
	Emoji       string    `json:"emoji"`
}

// 📐 Ингредиент в калькуляции
type CalculationIngredient struct {
	ID              string    `json:"id"`
	ProductName     string    `json:"product_name"`
	AmountGrams     int       `json:"amount_grams"`
	PricePerKg      float64   `json:"price_per_kg"`
	WastePercent    float64   `json:"waste_percent"`
	PriceAfterWaste float64   `json:"price_after_waste"`
	TotalCost       float64   `json:"total_cost"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

// 🧾 Финальная калькуляция блюда
type MenuCalculation struct {
	ID               string                  `json:"id"`
	MenuItemID       string                  `json:"menu_item_id"`
	TotalOutputGrams int                     `json:"total_output_grams"`
	TotalCost        float64                 `json:"total_cost"`
	Ingredients      []CalculationIngredient `json:"ingredients"`
	CreatedAt        time.Time               `json:"created_at"`
}

// ✅ Псевдоним для совместимости с utils.SaveDishCalculation
type Calculation = MenuCalculation



