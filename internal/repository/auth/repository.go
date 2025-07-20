package auth

import (
	"context"
	"database/sql"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/evgeniySeleznev/auth-project/internal/model"
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

func (r *repo) Create(ctx context.Context, info *model.User) (int64, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	//// Генерация фейковых данных
	//name := gofakeit.Name()
	//email := gofakeit.Email()
	//password := gofakeit.Password(true, true, true, true, true, 12)
	//role := 1

	//Вставка записи в таблицу users
	builderInsert := psql.Insert(tableName).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(info.Name, info.Email, info.Password, info.Role).
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

	log.Printf("Id: %d, Name: %s, email: %s, pass: %s, role: %v", userID, info.Name, info.Email, info.Password, info.Role)

	return int64(userID), nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Получение одной обновленной записи
	builderSelectOne := psql.Select(nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{"id": id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build select-one query: %v", err)
	}

	var role int
	var name, email, password string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = r.db.QueryRow(ctx, query, args...).Scan(&name, &email, &password, &role, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select updated user: %v", err)
	}

	log.Printf("Selected updated user -> name: %s, email: %s, password: %s, role: %d, created_at: %v, updated_at: %v",
		name, email, password, role, createdAt, updatedAt)

	return &model.User{
		Name:      name,
		Email:     email,
		Password:  password,
		Role:      model.Role(role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// Update ...
func (r *repo) Update(ctx context.Context, info *model.User) (*emptypb.Empty, error) {
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
func (r *repo) Delete(ctx context.Context, req *model.User) (*emptypb.Empty, error) {
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
