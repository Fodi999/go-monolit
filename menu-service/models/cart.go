package models

import "gorm.io/gorm"

// 🛒 Корзина пользователя
type Cart struct {
	gorm.Model
	UserID string     `gorm:"not null;index" json:"userId"` // индекс по UserID
	Items  []CartItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
}

// 📦 Позиция в корзине
// 📦 Позиция в корзине
type CartItem struct {
    gorm.Model
    CartID     uint    `gorm:"not null;index" json:"CartID"`   // теперь отдаём CartID
    MenuItemID string  `gorm:"not null"      json:"menuItemId"`
    Name       string  `gorm:"not null"      json:"name"`
    Quantity   int     `gorm:"not null"      json:"quantity"`
    Price      float64 `gorm:"not null"      json:"price"`
	ImageURL    string  `gorm:"-" json:"imageUrl"` // <— новое поле, не сохраняется в cart_items
}


// 🧾 Заказ пользователя
type Order struct {
	gorm.Model
	UserID     string      `gorm:"not null;index" json:"userId"`
	CartID     uint        `gorm:"index" json:"cartId"`
	Items      []OrderItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`
	TotalPrice float64     `gorm:"not null" json:"totalPrice"`
	Status     string      `gorm:"type:varchar(50);default:'pending'" json:"status"`
}

// 🧾 Позиция в заказе
type OrderItem struct {
	gorm.Model
	OrderID    uint    `gorm:"not null;index" json:"orderId"`
	MenuItemID string  `gorm:"not null" json:"menuItemId"`
	Name       string  `gorm:"not null" json:"name"`
	Quantity   int     `gorm:"not null" json:"quantity"`
	Price      float64 `gorm:"not null" json:"price"`
}

