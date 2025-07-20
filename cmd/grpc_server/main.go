package main

import (
	"context"
	"flag"
	"github.com/evgeniySeleznev/auth-project/internal/config"
	"github.com/evgeniySeleznev/auth-project/internal/config/env"
	"github.com/evgeniySeleznev/auth-project/internal/converter"
	authRepo "github.com/evgeniySeleznev/auth-project/internal/repository/auth"
	"github.com/evgeniySeleznev/auth-project/internal/service"
	authService "github.com/evgeniySeleznev/auth-project/internal/service/auth"
	desc "github.com/evgeniySeleznev/auth-project/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

// Create ...
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.authService.Create(ctx, converter.ToModelFromDesc(req))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted note with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

// Get ...
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("name: %s, email: %s, password: %s, role: %s,created_at: %v, updated_at: %v\n", user.Name, user.Email, user.Password, user.Role)

	return &desc.GetResponse{
		Name:  user.Name,
		Email: user.Email,
		Role:  desc.Role(user.Role),
	}, nil
}

// Update ...
//func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
//	_, err := s.authService.Update(ctx, &desc.User{})
//	if err != nil {
//		return nil, err
//	}
//
//	log.Println("update done")
//
//	return &emptypb.Empty{}, nil
//}

// Delete ...
//func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
//	_, err := s.authRepository.Delete(ctx, &desc.User{})
//	if err != nil {
//		return nil, err
//	}
//
//	log.Println("delete done")
//
//	return &emptypb.Empty{}, nil
//}

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

	authRepos := authRepo.NewRepository(pool)
	authSrv := authService.NewService(authRepos)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{authService: authSrv})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
