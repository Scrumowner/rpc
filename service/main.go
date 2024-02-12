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
	"net/rpc"
	"os"
	controller "rpc/service/controller"
	"rpc/service/internal/migrator"
	"rpc/service/models"
	pb "rpc/service/proto/gen"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can't load config from env")
	}
	configRpc := os.Getenv("RPC_PROVIDER")
	time.Sleep(10 * time.Second)
	dbConnection := fmt.Sprintf("user=user password=password host=db port=5432 dbname=my_database sslmode=disable")
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Fatal("Can't connect to PostgreSQL:", err)
	}
	defer db.Close()
	log.Println("Connect to postgres")

	dbx := sqlx.NewDb(db, "postgres")
	if err = dbx.Ping(); err != nil {
		log.Fatal("Can't use sqlx driver for PostgreSQL:", err)
	}
	address := models.SearchIntoDb{}
	geo := models.GeoIntoDb{}
	if err == nil {
		log.Println("Sqlx connect to postgres")
	}

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
	if configRpc == "json-rpc" {
		controller := controller.NewGeoControllerJsonRpc(client, redis, dbx)
		rpc.Register(controller)
		rpc.HandleHTTP()

		listener, err := net.Listen("tcp", "0.0.0.0:1234")
		if err != nil {
			log.Fatal("Can't open connection")
		}

		http.Serve(listener, nil)
	} else if configRpc == "rpc" {
		controller := controller.NewGeoController(client, redis, dbx)
		rpc.Register(controller)
		listener, err := net.Listen("tcp", "0.0.0.0:1234")
		if err != nil {
			log.Fatal("Can't open connection")
		}

		for {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				log.Println("Can't accept connection", acceptErr)
			}
			go rpc.ServeConn(conn)
		}
	} else if configRpc == "grpc" {
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
	}

	log.Println("Server is shut down")
}

//protoc -I . geo.proto --go_out=./gen/ --go_opt=paths=source_relative --go-grpc_out=./gen/ --go-grpc_opt=paths=source_relative
