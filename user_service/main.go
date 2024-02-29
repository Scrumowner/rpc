package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
	"user/controller"
	"user/models"
	pb "user/proto/user"
	"user/storage"
)

type DbConfig struct {
	user     string
	password string
	host     string
	port     string
	dbname   string
}
type ServConfig struct {
	listen string
	port   string
}

func main() {
	godotenv.Load()
	dbconfig := DbConfig{
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		host:     os.Getenv("POSTGRES_HOST"),
		port:     os.Getenv("POSTGRES_PORT"),
		dbname:   os.Getenv("POSTGRES_DB"),
	}
	servconfig := ServConfig{
		listen: os.Getenv("USER_LISTEN"),
		port:   os.Getenv("USER_PORT"),
	}
	time.Sleep(time.Second * 15)
	dbAddr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbconfig.user,
		dbconfig.password,
		dbconfig.host,
		dbconfig.port,
		dbconfig.dbname,
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
	migrator := storage.NewMigrator(dbx)
	user := &models.User{}
	migrator.Migrate(user)

	controller := controller.NewUserController(dbx)
	port := fmt.Sprintf("%s:%s", servconfig.listen, servconfig.port)
	conn, err := net.Listen("tcp", port)
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
