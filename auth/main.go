package main

import (
	controller "auth/controller"
	pb "auth/proto/auth"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type ServConfig struct {
	port string
}
type ConnConfig struct {
	host string
	port string
}

func main() {
	godotenv.Load()
	connconfig := &ConnConfig{
		host: os.Getenv("USER_HOST"),
		port: os.Getenv("USER_PORT"),
	}
	dial, err := grpc.Dial(fmt.Sprintf("%s:%s", connconfig.host, connconfig.port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	servconfig := &ServConfig{
		port: os.Getenv("AUTH_PORT"),
	}
	conn, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", servconfig.port))
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта: %v", err)
	}
	controller := controller.NewAuthController(dial)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, controller)
	log.Println("Starting server auth")
	if err = server.Serve(conn); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
