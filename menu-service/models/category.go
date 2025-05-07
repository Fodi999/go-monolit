package models

import "time"

type Category struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"` // 👉 теперь поле уникальное
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}



