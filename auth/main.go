package main

import (
	"auth/config"
	controller "auth/controller"
	pba "auth/proto/auth"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can't load config from .env files")
	}
	cfg := config.NewConfig()
	cfg.Load()
	userDial, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.UserService.Host, cfg.UserService.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen, cfg.Port))
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта: %v", err)
	}
	controller := controller.NewAuthController(userDial, cfg)

	server := grpc.NewServer()
	pba.RegisterAuthServiceServer(server, controller)
	log.Println("Starting server auth")
	if err = server.Serve(conn); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

//protoc -I . auth.proto --go_out=./proto/auth --go_opt=paths=source_relative --go-grpc_out=./proto/auth --go-grpc_opt=paths=source_relative
//export PATH="$PATH:$(go env GOPATH)/bin"
