package handlers

import (
	"context"
	"log"
	"net/mail"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"

	"monolith/auth_service/database"
	"monolith/auth_service/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Role    string `json:"role"`
	ID      string `json:"id"`
	Initial string `json:"initial"` // 💡 первая буква имени
}

type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Bio      string `json:"bio"`
	Birthday string `json:"birthday"` // Формат: "2006-01-02"
}

func HandleRegister(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный формат запроса")
	}

	var exists bool
	err := database.DB.QueryRow(context.Background(), `
		SELECT EXISTS(
			SELECT 1 FROM users 
			WHERE email=$1 OR phone=$2 OR name=$3
		)
	`, req.Email, req.Phone, req.Name).Scan(&exists)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при проверке уникальности")
	}
	if exists {
		return fiber.NewError(fiber.StatusConflict, "Email, имя или телефон уже используется")
	}

	id := uuid.NewString()
	_, err = database.DB.Exec(context.Background(),
		"INSERT INTO users (id, email, password, role, name, phone, created_at) VALUES ($1, $2, $3, $4, $5, $6, NOW())",
		id, req.Email, req.Password, "user", req.Name, req.Phone)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при создании пользователя")
	}

	token, err := utils.GenerateToken(id, "user")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка генерации токена")
	}

	initial := ""
	if len(req.Name) > 0 {
		initial = string([]rune(req.Name)[0])
	}

	return c.JSON(LoginResponse{Token: token, Role: "user", ID: id, Initial: initial})
}

func HandleLogin(c *fiber.Ctx) error {
	var creds LoginRequest
	if err := c.BodyParser(&creds); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный формат запроса")
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminID := os.Getenv("ADMIN_ID")

	if creds.Email == adminEmail && creds.Password == adminPassword {
		token, err := utils.GenerateToken(adminID, "admin")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Ошибка генерации токена")
		}
		log.Printf("🛡️ Вход админа: %s\n", adminID)
		return c.JSON(LoginResponse{Token: token, Role: "admin", ID: adminID, Initial: "A"})
	}

	var id, password, role, name string
	err := database.DB.QueryRow(context.Background(),
		"SELECT id, password, role, name FROM users WHERE email=$1", creds.Email).
		Scan(&id, &password, &role, &name)
	if err != nil || password != creds.Password {
		log.Printf("❌ Ошибка входа: пользователь не найден или неверный пароль для %s\n", creds.Email)
		return fiber.NewError(fiber.StatusUnauthorized, "Неверный email или пароль")
	}

	token, err := utils.GenerateToken(id, role)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка генерации токена")
	}

	initial := ""
	if len(name) > 0 {
		initial = string([]rune(name)[0])
	}

	log.Printf("✅ Вход пользователя: %s (%s)\n", id, creds.Email)
	_, err = database.DB.Exec(context.Background(), "UPDATE users SET last_active = NOW() WHERE id=$1", id)
	if err != nil {
		log.Printf("❌ Ошибка обновления last_active: %v\n", err)
	} else {
		log.Printf("🕒 last_active успешно обновлён для %s\n", id)
	}

	return c.JSON(LoginResponse{Token: token, Role: role, ID: id, Initial: initial})
}

func HandleCheckEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email не указан")
	}

	var exists bool
	err := database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при проверке email")
	}

	return c.JSON(fiber.Map{"exists": exists})
}

func HandleCheckPhone(c *fiber.Ctx) error {
	phone := c.Query("phone")
	if phone == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Телефон не указан")
	}

	var exists bool
	err := database.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE phone=$1)", phone).Scan(&exists)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при проверке телефона")
	}

	return c.JSON(fiber.Map{"exists": exists})
}

func HandleUpdateProfile(c *fiber.Ctx) error {
	claims := c.Locals("claims").(jwt.MapClaims)
	userID := claims["id"].(string)

	var req UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Неверный формат запроса")
	}

	if req.Name == "" || req.Email == "" || req.Phone == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Имя, email и телефон обязательны")
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Некорректный email")
	}

	var exists bool
	err := database.DB.QueryRow(context.Background(), `
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE (email=$1 OR phone=$2) AND id != $3
		)
	`, req.Email, req.Phone, userID).Scan(&exists)
	if err != nil {
		log.Printf("❌ Ошибка при проверке уникальности: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка при проверке данных")
	}
	if exists {
		return fiber.NewError(fiber.StatusConflict, "Email или телефон уже используется")
	}

	var birthday *time.Time
	if req.Birthday != "" {
		parsed, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Неверный формат даты (ожидается YYYY-MM-DD)")
		}
		birthday = &parsed
	}

	_, err = database.DB.Exec(context.Background(), `
		UPDATE users
		SET name=$1, email=$2, phone=$3, address=$4, bio=$5, birthday=$6
		WHERE id=$7
	`, req.Name, req.Email, req.Phone, req.Address, req.Bio, birthday, userID)
	if err != nil {
		log.Printf("❌ Ошибка при обновлении профиля: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось обновить профиль")
	}

	log.Printf("✅ Профиль обновлён: %s", userID)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Профиль успешно обновлён",
	})
}




