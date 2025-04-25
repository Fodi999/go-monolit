package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Phone      string  `json:"phone"`
	Role       string  `json:"role"`
	Password   string  `json:"password"`
	Address    *string `json:"address,omitempty"`
	Bio        *string `json:"bio,omitempty"`
	Birthday   *string `json:"birthday,omitempty"`
	CreatedAt  string  `json:"created_at"`
	LastActive string  `json:"last_active"`
	Avatar     string  `json:"avatar"`
}

// ✅ Получить всех пользователей
func HandleGetAllUsers(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Println("📥 Запрос на получение всех пользователей")

		claims := c.Locals("claims").(jwt.MapClaims)
		if claims["role"] != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "Доступ запрещен")
		}

		rows, err := db.Query(context.Background(), `
			SELECT 
				id, name, email, phone, role, password,
				address, bio,
				TO_CHAR(birthday, 'YYYY-MM-DD') as birthday,
				TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
				TO_CHAR(last_active, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as last_active
			FROM users
			ORDER BY created_at DESC
		`)
		if err != nil {
			log.Printf("❌ Ошибка при запросе пользователей: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при получении пользователей")
		}
		defer rows.Close()

		var users []UserResponse
		for rows.Next() {
			var u UserResponse
			err := rows.Scan(
				&u.ID,
				&u.Name,
				&u.Email,
				&u.Phone,
				&u.Role,
				&u.Password,
				&u.Address,
				&u.Bio,
				&u.Birthday,
				&u.CreatedAt,
				&u.LastActive,
			)
			if err != nil {
				log.Printf("❌ Ошибка при обработке строки: %v", err)
				return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при обработке данных")
			}
			u.Password = ""
			if len(u.Name) > 0 {
				runes := []rune(u.Name)
				u.Avatar = string(runes[0])
			} else {
				u.Avatar = "?"
			}
			users = append(users, u)
		}

		log.Printf("✅ Получено пользователей: %d", len(users))
		return c.JSON(users)
	}
}

// ✅ Удалить пользователя
func HandleDeleteUser(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		log.Printf("🗑️ Запрос на удаление пользователя: %s", userID)

		claims := c.Locals("claims").(jwt.MapClaims)
		if claims["role"] != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "Доступ запрещен")
		}

		cmdTag, err := db.Exec(context.Background(), "DELETE FROM users WHERE id=$1", userID)
		if err != nil {
			log.Printf("❌ Ошибка при удалении пользователя: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при удалении пользователя")
		}

		if cmdTag.RowsAffected() == 0 {
			return fiber.NewError(fiber.StatusNotFound, "Пользователь не найден")
		}

		log.Printf("✅ Пользователь %s удалён", userID)
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// ✅ Получить пользователя по ID
func HandleGetUserByID(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		log.Printf("📥 Запрос на получение профиля пользователя: %s", userID)

		var user UserResponse
		err := db.QueryRow(context.Background(), `
			SELECT 
				id, name, email, phone, role, password,
				address, bio,
				TO_CHAR(birthday, 'YYYY-MM-DD') as birthday,
				TO_CHAR(created_at, 'YYYY-MM-DD"T"HH24:MI:SS') as created_at,
				TO_CHAR(last_active, 'YYYY-MM-DD"T"HH24:MI:SS') as last_active
			FROM users
			WHERE id = $1
		`, userID).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Role,
			&user.Password,
			&user.Address,
			&user.Bio,
			&user.Birthday,
			&user.CreatedAt,
			&user.LastActive,
		)

		if err != nil {
			log.Printf("❌ Ошибка при получении пользователя: %v", err)
			return fiber.NewError(fiber.StatusNotFound, "Пользователь не найден")
		}

		user.Password = ""
		if len(user.Name) > 0 {
			runes := []rune(user.Name)
			user.Avatar = string(runes[0])
		} else {
			user.Avatar = "?"
		}

		return c.JSON(user)
	}
}

// ✅ Обновить роль пользователя (admin может назначить любую разрешённую роль)
func HandleUpdateUserRole(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(jwt.MapClaims)
		if claims["role"] != "admin" {
			return fiber.NewError(fiber.StatusForbidden, "Доступ запрещён")
		}

		userID := c.Params("id")
		var body struct {
			Role string `json:"role"`
		}
		if err := c.BodyParser(&body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Неверный формат тела запроса")
		}

		// ✅ Список разрешённых ролей
		allowedRoles := map[string]bool{
			"user":     true,
			"admin":    true,
			"повар":    true,
			"курьер":   true,
			"официант": true,
		}
		if !allowedRoles[body.Role] {
			return fiber.NewError(fiber.StatusBadRequest, "Недопустимая роль")
		}

		_, err := db.Exec(context.Background(), `UPDATE users SET role=$1 WHERE id=$2`, body.Role, userID)
		if err != nil {
			log.Printf("❌ Ошибка обновления роли пользователя %s: %v", userID, err)
			return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при обновлении роли")
		}

		log.Printf("✅ Роль пользователя %s обновлена на %s", userID, body.Role)
		return c.JSON(fiber.Map{
			"success": true,
			"role":    body.Role,
		})
	}
}





