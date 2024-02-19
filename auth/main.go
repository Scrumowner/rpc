package main

import (
	controller "auth/controller"
	pb "auth/proto/auth"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	dial, err := grpc.Dial("user_service:1237", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	conn, err := net.Listen("tcp", "0.0.0.0:1235")
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
