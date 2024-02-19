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
	"user/controller"
	"user/models"
	pb "user/proto/user"
	"user/storage"
)

func main() {
	time.Sleep(time.Second * 15)
	dbAddr := fmt.Sprintf("user=user password=password host=db port=5432 dbname=my_database sslmode=disable")
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
	migrator := storage.NewMigrator(dbx)
	user := &models.User{}
	migrator.Migrate(user)

	controller := controller.NewUserController(dbx)
	conn, err := net.Listen("tcp", ":1237")
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, controller)
	log.Println("Starting server auth")
	err = server.Serve(conn)
	if err != nil {
		log.Println("Can't serve grpc user_service")
	}

}
