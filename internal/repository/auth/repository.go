package auth

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit"
	desc "github.com/evgeniySeleznev/auth-project/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/evgeniySeleznev/auth-project/internal/repository"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool //пул коннектов к постгресу
}

func NewRepository(db *pgxpool.Pool) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *desc.User) (int64, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Генерация фейковых данных
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, true, 12)
	role := 1

	//Вставка записи в таблицу users
	builderInsert := psql.Insert(tableName).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(name, email, password, role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build insert query: %v", err)
	}

	var userID int
	err = r.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("Name: %s, email: %s, pass: %s, role: %v", info.GetName(), info.GetEmail(), info.GetPassword(), info.GetRole())

	return int64(userID), nil
}

func (r *repo) Get(ctx context.Context, id int64) (*desc.User, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Получение одной обновленной записи
	builderSelectOne := psql.Select(idColumn, nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build select-one query: %v", err)
	}

	var idz, re int
	var n, e, p string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = r.db.QueryRow(ctx, query, args...).Scan(&idz, &n, &e, &p, &r, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select updated user: %v", err)
	}

	log.Printf("Selected updated user -> id: %d, name: %s, email: %s, password: %s, role: %d, created_at: %v, updated_at: %v",
		id, n, e, p, re, createdAt, updatedAt)

	return &desc.User{
		Name:  n,
		Email: e,
		Role:  desc.Role(re),
	}, nil
}

// Update ...
func (r *repo) Update(ctx context.Context, info *desc.User) (*emptypb.Empty, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builderUpdate := psql.Update("users").
		Set("name", gofakeit.Name()).
		Set("email", gofakeit.Email()).
		Set("password", gofakeit.Password(true, true, true, true, true, 14)).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": gofakeit.UUID()})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build update query: %v", err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to update user: %v", err)
	}

	log.Printf("Updated %d rows", res.RowsAffected())

	return &emptypb.Empty{}, nil
}

// Delete ...
func (r *repo) Delete(ctx context.Context, req *desc.User) (*emptypb.Empty, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	// Удаление записи
	builderDelete := psql.Delete("users").
		Where(sq.Eq{"id": gofakeit.UUID()})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Fatalf("failed to build delete query: %v", err)
	}

	res, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		log.Fatalf("failed to delete user: %v", err)
	}

	log.Printf("Deleted %d rows", res.RowsAffected())

	log.Printf("Delete qequest ID: %d", gofakeit.UUID())
	return &emptypb.Empty{}, nil
}
