package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
	"net/http"
	"net/rpc"
	controller "rpc/service/controller"
	"rpc/service/internal/migrator"
	"rpc/service/models"
	"time"
)

// service using rpc

func main() {
	time.Sleep(10 * time.Second)
	// Формирование строки подключения к PostgreSQL
	dbConnection := fmt.Sprintf("user=user password=password host=db port=5432 dbname=my_database sslmode=disable")

	// Подключение к базе данных PostgreSQL
	db, err := sql.Open("postgres", dbConnection)
	if err != nil {
		log.Fatal("Can't connect to PostgreSQL:", err)
	}
	defer db.Close()
	log.Println("Connect to postgres")

	// Создание объекта sqlx.DB для удобной работы с базой данных
	dbx := sqlx.NewDb(db, "postgres")
	if err = dbx.Ping(); err != nil {
		log.Fatal("Can't use sqlx driver for PostgreSQL:", err)
	}
	address := models.SearchIntoDb{}
	geo := models.GeoIntoDb{}
	log.Println(err)
	log.Println("Sqlx connect to postgres")
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
}
