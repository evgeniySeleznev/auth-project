package main

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
)

func main() {
	ctx := context.Background()

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Генерация фейковых данных
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, 12)
	role := 1

	//Вставка записи в таблицу users
	builderInsert := psql.Insert("users").
		Columns("name", "email", "password", "role").
		Values(name, email, password, role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build insert query: %v", err)
	}

	var userID int
	err = pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("Inserted user with id: %d", userID)

	// Выборка всех пользователей
	builderSelect := psql.Select("id", "name", "email", "password", "role", "created_at", "updated_at").
		From("users").
		OrderBy("id ASC").
		Limit(10)

	query, args, err = builderSelect.ToSql()
	if err != nil {
		log.Fatalf("failed to build select query: %v", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to select users: %v", err)
	}
	defer rows.Close()

	var id, r int
	var n, e, p string
	var createdAt time.Time
	var updatedAt sql.NullTime

	for rows.Next() {
		err = rows.Scan(&id, &n, &e, &p, &r, &createdAt, &updatedAt)
		if err != nil {
			log.Fatalf("failed to scan user: %v", err)
		}

		log.Printf("id: %d, name: %s, email: %s, password: %s, role: %d, created_at: %v, updated_at: %v",
			id, n, e, p, r, createdAt, updatedAt)
	}

	// Обновление записи
	builderUpdate := psql.Update("users").
		Set("name", gofakeit.Name()).
		Set("email", gofakeit.Email()).
		Set("password", gofakeit.Password(true, true, true, true, true, 14)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": userID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build update query: %v", err)
	}

	res, err := pool.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf("Updated %d rows", res.RowsAffected())

	// Получение одной обновленной записи
	builderSelectOne := psql.Select("id", "name", "email", "password", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": userID}).
		Limit(1)

	query, args, err = builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build select-one query: %v", err)
	}

	err = pool.QueryRow(ctx, query, args...).Scan(&id, &n, &e, &p, &r, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select updated user: %v", err)
	}

	log.Printf("Selected updated user -> id: %d, name: %s, email: %s, password: %s, role: %d, created_at: %v, updated_at: %v",
		id, n, e, p, r, createdAt, updatedAt)
}
