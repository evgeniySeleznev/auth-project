package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
)

const (
	dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	// Делаем запрос на вставку записи в таблицу note
	res, err := con.Exec(ctx, "INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)", gofakeit.Name(), gofakeit.Email(), gofakeit.Password(true, true, true, true, true, 12), 1)
	if err != nil {
		log.Fatalf("failed to insert note: %v", err)
	}

	log.Printf("inserted %d rows", res.RowsAffected())

	// Делаем запрос на выборку записей из таблицы note
	rows, err := con.Query(ctx, "SELECT id, name, email, password, role, created_at, updated_at FROM users")
	if err != nil {
		log.Fatalf("failed to select notes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, role int
		var name, email, password string
		var createdAt time.Time
		var updatedAt sql.NullTime

		err = rows.Scan(&id, &name, &email, &password, &role, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan note: %v", err)
		}

		log.Printf("id: %d, name: %s, email: %s, password: %s, role: %d, created_at: %v, updated_at: %v\n", id, name, email, password, role, createdAt, updatedAt)
	}
}
