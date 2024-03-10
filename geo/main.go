package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"rpc/service/config"
	controller "rpc/service/controller"
	"rpc/service/internal/migrator"
	"rpc/service/models"
	pb "rpc/service/proto/geo"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can't load .env geo service")
	}
	time.Sleep(10 * time.Second)
	cfg := config.NewConfig()
	cfg.Load()
	dbAddr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DbConfig.User,
		cfg.DbConfig.Password,
		cfg.DbConfig.Host,
		cfg.DbConfig.Port,
		cfg.DbConfig.Dbname,
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

	address := models.SearchIntoDb{}
	geo := models.GeoIntoDb{}

	migrator := migrator.NewMigrator(dbx)
	migrator.Migrate(&address, &geo)
	cache := fmt.Sprintf("%s:%s", cfg.CacheConfig.Host, cfg.CacheConfig.Port)
	redis := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         cache,
		DB:           0,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 20,
	})
	client := http.Client{}

	controller := controller.NewGeoContollergRpc(client, redis, dbx, cfg)
	port := fmt.Sprintf("%s:%s", cfg.Listen, cfg.Port)
	listner, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Can't open connection")
	}
	server := grpc.NewServer()
	pb.RegisterGeoServiceServer(server, controller)
	log.Println("grpc server starting")
	if err := server.Serve(listner); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	log.Println("Server is shut down")
}

//protoc -I . geo.proto --go_out=./proto/geo --go_opt=paths=source_relative --go-grpc_out=./proto/geo --go-grpc_opt=paths=source_relative --grpc-gateway_out=./proto/geo --grpc-gateway_opt=paths=source_relative --openapiv2_out=./proto/geo
//protoc -I . geo.proto --go_out=./proto/geo --go_opt=paths=source_relative --go-grpc_out=./proto/geo --go-grpc_opt=paths=source_relative --grpc-gateway_out=./proto/geo --grpc-gateway_opt=paths=source_relative --openapiv2_out=./proto/geo
//$ export PATH="$PATH:$(go env GOPATH)/bin"
