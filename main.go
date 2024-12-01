package main

import (
	"http-server/danilkovalev/internal/repository"
    "http-server/danilkovalev/internal/service"
	pb "http-server/danilkovalev/internal/proto"
	grpcHandler "http-server/danilkovalev/internal/transport/grpc"
    "http-server/danilkovalev/internal/transport/rest"

    "log"
	"net"
    "github.com/joho/godotenv"
	"google.golang.org/grpc"
)


func main() {
    err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mySqlDB, err := repository.NewMySQLDB()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL database: %v", err)
	}

	err = repository.Migrations(mySqlDB)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	beanstalkClient, err := service.NewBeanstalkClient()
	if err != nil {
		log.Fatalf("Error creating Beanstalk client: %v", err)
	}

    repository := repository.NewRepository(mySqlDB)
	service := service.NewService(repository, beanstalkClient)
	httphandler := handler.NewHandler(service)
	grpcServer := grpcHandler.NewGRPCServer(service)

    go func() {
		httphandler.InitRoutes()
	}()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, grpcServer)

	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
