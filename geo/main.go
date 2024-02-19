package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	controller "rpc/service/controller"
	"rpc/service/internal/migrator"
	"rpc/service/models"
	pb "rpc/service/proto/geo"
	"time"
)

func main() {
	time.Sleep(10 * time.Second)
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

	address := models.SearchIntoDb{}
	geo := models.GeoIntoDb{}

	migrator := migrator.NewMigrator(dbx)
	migrator.Migrate(&address, &geo)

	redis := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         "cache:6379",
		DB:           0,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 20,
	})
	client := http.Client{}

	controller := controller.NewGeoContollergRpc(client, redis, dbx)
	listner, err := net.Listen("tcp", "0.0.0.0:1234")
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
