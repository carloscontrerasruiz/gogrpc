package main

import (
	"fmt"
	"log"
	"net"

	"github.com/carloscontrerasruiz/gogrpc/database"
	"github.com/carloscontrerasruiz/gogrpc/server"
	"github.com/carloscontrerasruiz/gogrpc/testspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", "localhost:5070")
	if err != nil {
		log.Fatal(err)
	}

	dbName := "postgres"
	url := fmt.Sprintf("postgres://postgres:postgres@localhost:5432/%s?sslmode=disable", dbName)

	repo, err := database.NewPostgresRepository(url)
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewTestServer(repo)
	s := grpc.NewServer()
	testspb.RegisterTestServiceServer(s, server)

	reflection.Register(s)

	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
