package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Avatar     string  `json:"avatar"` // первая буква имени
	Email      string  `json:"email"`
	Phone      string  `json:"phone"`
	Role       string  `json:"role"`
	Address    *string `json:"address,omitempty"`
	Bio        *string `json:"bio,omitempty"`
	Birthday   *string `json:"birthday,omitempty"`
	CreatedAt  string  `json:"created_at"`
	LastActive string  `json:"last_active"`
	Online     bool    `json:"online"`
	Orders     int     `json:"orders"`
}

func HandleProfile(c *fiber.Ctx) error {
	claims := c.Locals("claims").(jwt.MapClaims)
	userID := claims["id"].(string)

	db := c.Locals("db").(*pgxpool.Pool)

	var p ProfileResponse
	err := db.QueryRow(context.Background(), `
		SELECT 
			id, name, email, phone, role,
			address, bio,
			TO_CHAR(birthday, 'YYYY-MM-DD') as birthday,
			TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS') as created_at,
			TO_CHAR(last_active, 'YYYY-MM-DD"T"HH24:MI:SS') as last_active,
			online, orders
		FROM users
		WHERE id = $1
	`, userID).Scan(
		&p.ID,
		&p.Name,
		&p.Email,
		&p.Phone,
		&p.Role,
		&p.Address,
		&p.Bio,
		&p.Birthday,
		&p.CreatedAt,
		&p.LastActive,
		&p.Online,
		&p.Orders,
	)
	if err != nil {
		log.Printf("❌ Ошибка при загрузке профиля: %v", err)
		return fiber.NewError(fiber.StatusNotFound, "Профиль не найден")
	}

	if len(p.Name) > 0 {
		p.Avatar = string([]rune(p.Name)[0])
	} else {
		p.Avatar = "?"
	}

	return c.JSON(p)
}
