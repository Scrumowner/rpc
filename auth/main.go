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
	listen string
	port   string
}
type ConnConfig struct {
	host string
	port string
}

func main() {
	godotenv.Load(".env", "auth.env")
	connconfig := &ConnConfig{
		host: os.Getenv("USER_HOST"),
		port: os.Getenv("USER_PORT"),
	}
	dial, err := grpc.Dial(fmt.Sprintf("%s:%s", connconfig.host, connconfig.port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	servconfig := &ServConfig{
		listen: os.Getenv("AUTH_LISTEN"),
		port:   os.Getenv("AUTH_PORT"),
	}
	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%s", servconfig.listen, servconfig.port))
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
