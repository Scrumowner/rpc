package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
	"user/config"
	"user/internal/infra/migrator"
	"user/internal/models"
	"user/internal/modules"
	pb "user/proto"
)

func main() {

	cfg := config.NewcConfig()
	cfg.Load()
	time.Sleep(time.Second * 15)
	dbAddr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
	)
	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		log.Fatal("Can't connect to PostgreSQL:", err)
	}
	defer db.Close()
	log.Println("Connect to postgres")

	dbx := sqlx.NewDb(db, "postgres")
	if err = dbx.Ping(); err != nil {
		log.Fatal("Can't use sqlx driver for PostgreSQL:", err)
	}
	if err == nil {
		log.Println("Sqlx connect to postgres")
	}
	migrator := migrator.NewMigrator(dbx)
	user := &models.User{}
	migrator.Migrate(user)

	controller := modules.NewControllers(dbx)

	port := fmt.Sprintf("%s:%s", cfg.ServConfig.Listen, cfg.ServConfig.Port)

	conn, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, controller.User)
	log.Println("Starting server user_service")
	err = server.Serve(conn)
	if err != nil {
		log.Println("Can't serve grpc user_service")
	}

}

//protoc -I . user.proto --go_out=./proto/ --go_opt=paths=source_relative --go-grpc_out=./proto/ --go-grpc_opt=paths=source_relative
//export PATH="$PATH:$(go env GOPATH)/bin"
