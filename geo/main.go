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
	"os"
	controller "rpc/service/controller"
	"rpc/service/internal/migrator"
	"rpc/service/models"
	pb "rpc/service/proto/geo"
	"time"
)

type DbConfig struct {
	user     string
	password string
	host     string
	port     string
	dbname   string
}
type CacheConfig struct {
	host string
	port string
}
type ServConfig struct {
	listen string
	port   string
}

func main() {
	godotenv.Load(".env", "service.env")
	dbconfig := &DbConfig{
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		host:     os.Getenv("POSTGRES_HOST"),
		port:     os.Getenv("POSTGRES_PORT"),
		dbname:   os.Getenv("POSTGRES_DB"),
	}
	cacheconfig := &CacheConfig{
		host: os.Getenv("REDIS_HOST"),
		port: os.Getenv("REDIS_PORT"),
	}
	servconfig := &ServConfig{
		listen: os.Getenv("USER_LISTEN"),
		port:   os.Getenv("USER_PORT"),
	}
	time.Sleep(10 * time.Second)

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

	address := models.SearchIntoDb{}
	geo := models.GeoIntoDb{}

	migrator := migrator.NewMigrator(dbx)
	migrator.Migrate(&address, &geo)
	cache := fmt.Sprintf("%s:%s", cacheconfig.host, cacheconfig.port)
	redis := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         cache,
		DB:           0,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 20,
	})
	client := http.Client{}

	controller := controller.NewGeoContollergRpc(client, redis, dbx)
	port := fmt.Sprintf("%s:%s", servconfig.listen, servconfig.port)
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
