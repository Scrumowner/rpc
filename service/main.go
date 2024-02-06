package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log"
	"net"
	"net/http"
	"net/rpc"
	controller "rpc/service/controller"
	"time"
)

//service using rpc

func main() {
	client := http.Client{}
	db, err := sql.Open("postgres", "")
	if err != nil {
		log.Fatal("Can't connect to postgres")
	}
	dbx := sqlx.NewDb(db, "postgres")
	if err = dbx.Ping(); err != nil {
		log.Fatal("Can't use sqlx driver for postgres")
	}
	log.Println(err)
	redis := redis.NewClient(&redis.Options{
		Network:      "tcp",
		Addr:         "",
		Username:     "",
		Password:     "",
		DB:           0,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 20,
	})
	controller := controller.NewGeoController(client, redis, dbx)
	rpc.Register(controller)
	listener, err := net.Listen("tcp", "8081")
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
