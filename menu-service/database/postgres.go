package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitPostgres() {
	// ⚙️ Используем отдельную переменную для menu-сервиса
	connStr := os.Getenv("MENU_DATABASE_URL")
	if connStr == "" {
		log.Fatal("❌ MENU_DATABASE_URL не найден в .env")
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к базе данных: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ Ошибка пинга базы данных: %v", err)
	}

	fmt.Println("✅ Успешное подключение к базе данных MENU")

	// 🔌 Расширение pgcrypto
	_, err = DB.Exec(`DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pgcrypto') THEN
			CREATE EXTENSION pgcrypto;
		END IF;
	END $$;`)
	if err != nil {
		log.Fatalf("❌ Ошибка активации расширения pgcrypto: %v", err)
	}
	fmt.Println("✅ Расширение pgcrypto проверено или включено")

	// 📦 Только если в режиме разработки
	if os.Getenv("APP_ENV") == "development" {
		fmt.Println("⚠️ DEV MODE: пересоздание таблиц...")

		_, _ = DB.Exec(`DROP TABLE IF EXISTS calculation_ingredients CASCADE;`)
		_, _ = DB.Exec(`DROP TABLE IF EXISTS menu_calculations CASCADE;`)
		_, _ = DB.Exec(`DROP TABLE IF EXISTS inventory_items CASCADE;`)
		_, _ = DB.Exec(`DROP TABLE IF EXISTS menu_items CASCADE;`)
		_, _ = DB.Exec(`DROP TABLE IF EXISTS categories CASCADE;`)
	}

	// Создание таблиц
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS categories (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatalf("❌ Ошибка создания таблицы categories: %v", err)
	}
	fmt.Println("✅ Таблица categories проверена")

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS menu_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		description TEXT,
		price NUMERIC(10, 2),
		cost_price NUMERIC(10, 2),
		image_url TEXT,
		margin NUMERIC(10, 2),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
		published BOOLEAN NOT NULL DEFAULT FALSE
	);`)
	if err != nil {
		log.Fatalf("❌ Ошибка создания таблицы menu_items: %v", err)
	}
	fmt.Println("✅ Таблица menu_items проверена")

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS inventory_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		product_name TEXT NOT NULL,
		weight_grams INTEGER NOT NULL,
		price_per_kg NUMERIC(10, 2) NOT NULL,
		available BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		emoji TEXT DEFAULT '🍽️'
	);`)
	if err != nil {
		log.Fatalf("❌ Ошибка создания таблицы inventory_items: %v", err)
	}
	fmt.Println("✅ Таблица inventory_items проверена")

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS menu_calculations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		menu_item_id UUID NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
		total_output_grams INTEGER NOT NULL,
		total_cost NUMERIC(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatalf("❌ Ошибка создания таблицы menu_calculations: %v", err)
	}
	fmt.Println("✅ Таблица menu_calculations проверена")

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS calculation_ingredients (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		calculation_id UUID NOT NULL REFERENCES menu_calculations(id) ON DELETE CASCADE,
		product_name TEXT NOT NULL,
		amount_grams INTEGER NOT NULL,
		price_per_kg NUMERIC(10, 2) NOT NULL,
		waste_percent NUMERIC(5, 2) DEFAULT 0.0,
		price_after_waste NUMERIC(10, 2) NOT NULL,
		total_cost NUMERIC(10, 2) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatalf("❌ Ошибка создания таблицы calculation_ingredients: %v", err)
	}
	fmt.Println("✅ Таблица calculation_ingredients проверена")
}













