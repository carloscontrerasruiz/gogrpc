package main

import (
	"fmt"
	"log"
	"net"

	"github.com/carloscontrerasruiz/gogrpc/database"
	"github.com/carloscontrerasruiz/gogrpc/server"
	"github.com/carloscontrerasruiz/gogrpc/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", "localhost:5060")
	if err != nil {
		log.Fatal(err)
	}

	dbName := "postgres"
	url := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)

	repo, err := database.NewPostgresRepository(url)
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewStudentServer(repo)
	s := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
