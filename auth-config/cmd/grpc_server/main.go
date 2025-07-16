package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/evgeniySeleznev/auth-project/config/internal/config"
	"github.com/evgeniySeleznev/auth-project/config/internal/config/env"
	"log"
	"net"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/evgeniySeleznev/auth-project/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

// Create ...
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов

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
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}

	log.Printf("Name: %s, email: %s, pass: %s, pass_confirm: %s, role: %v", req.GetName(), req.GetEmail(), req.GetPassword(), req.GetPasswordConfirm(), req.GetRole())

	return &desc.CreateResponse{
		Id: int64(userID),
	}, nil
}

// Get ...
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Получение одной обновленной записи
	builderSelectOne := psql.Select("id", "name", "email", "password", "role", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		log.Fatalf("failed to build select-one query: %v", err)
	}

	var id, r int
	var n, e, p string
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = s.pool.QueryRow(ctx, query, args...).Scan(&id, &n, &e, &p, &r, &createdAt, &updatedAt)
	if err != nil {
		log.Fatalf("failed to select updated user: %v", err)
	}

	log.Printf("Selected updated user -> id: %d, name: %s, email: %s, password: %s, role: %d, created_at: %v, updated_at: %v",
		id, n, e, p, r, createdAt, updatedAt)

	var updatedAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtTime = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        int64(id),
		Name:      n,
		Email:     e,
		Role:      desc.Role(r),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: updatedAtTime,
	}, nil
}

// Update ...
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов
	log.Printf("Update request ID: %d, name: %s, email: %s", req.GetId(), req.GetName(), req.GetEmail())
	return &emptypb.Empty{}, nil
}

// Delete ...
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	_ = ctx // <– подавляем линтер, без логических побочных эффектов
	log.Printf("Delete qequest ID: %d", req.GetId())
	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	// Считываем переменные окружения
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
