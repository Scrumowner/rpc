package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"hugoproxy-main/proxy/controllers"
	internalMiddleware "hugoproxy-main/proxy/middleware"
	"hugoproxy-main/proxy/migrator"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"time"
)

const (
	host     = "db"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "my_database"
)

func main() {
	time.Sleep(time.Second * 5)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	dbx := sqlx.NewDb(db, "postgres")
	migrator := migrator.NewMigrator(dbx)
	migrator.Migrate()
	client := http.Client{}
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	cache := redis.NewClient(
		&redis.Options{
			Addr:     "cache:6379",
			Password: "",
			DB:       0,
		})
	pong, err := cache.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Can't connect to redis", err, pong)
	}
	fmt.Println("Successfully connected to Redis", pong)
	controller := controllers.NewControllers(sugar, client, dbx, cache)
	rp := internalMiddleware.NewReverseProxy("http://hugo", ":1313")
	r.Use(rp.ReverseProxy)
	r.Post("/api/register", controller.AuthController.Register)
	r.Post("/api/login", controller.AuthController.Login)
	r.Route("/swagger", func(r chi.Router) {
		r.Get("/index", controller.SwagController.GetSwaggerHtml)
		r.Get("/swagger", controller.SwagController.GetSwaggerJson)
	})
	r.Route("/mycustompath", func(r chi.Router) {
		//r.Use(jwtauth.Verifier(storage.TokenAuth))
		//r.Use(jwtauth.Authenticator(storage.TokenAuth))
		//r.Use(internalMiddleware.TokenAuthMiddleware)
		r.Get("/pprof/profile", pprof.Profile)
		r.Get("/pprof/trace", pprof.Trace)
		r.Get("/pprof/", pprof.Index)
		r.Get("/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
		r.Get("/pprof/block", pprof.Handler("block").ServeHTTP)
		r.Get("/pprof/cmdline", pprof.Cmdline)
		r.Get("/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
		r.Get("/pprof/heap", pprof.Handler("heap").ServeHTTP)
		r.Get("/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
		r.Get("/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		r.Get("/debug/pprof/goroutine", fullGoroutineStackDump)
	})

	r.Route("/api", func(r chi.Router) {
		//r.Use(jwtauth.Verifier(storage.TokenAuth))
		//r.Use(jwtauth.Authenticator(storage.TokenAuth))
		//r.Use(internalMiddleware.TokenAuthMiddleware)

		r.Route("/address", func(r chi.Router) {

			r.Post("/search", controller.SearchController.GetSearch)
			r.Post("/geocode", controller.SearchController.GetGeoCode)
		})
	})
	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	go func() {
		log.Printf(fmt.Sprintf("Server running on port %s", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("hutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
func fullGoroutineStackDump(w http.ResponseWriter, r *http.Request) {
	// Получаем полный стек горутин
	buf := make([]byte, 1<<20)
	stackLen := runtime.Stack(buf, true)

	// Отправляем содержимое стека в ответ на запрос
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(buf[:stackLen])
}

//{"id":123,"name":"123","phone":"123","email":"123","password":"123"}
