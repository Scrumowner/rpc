package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"proxy/internal/infra/ratelimit"
	"proxy/internal/router"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net/http"

	"proxy/config"
	"proxy/internal/modules"

	internalMiddleware "proxy/internal/middleware"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Proxy can't read .env")
	}
	conn := config.NewConfig()
	conn.Load()
	time.Sleep(time.Second * 15)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	geoConn, err := grpc.Dial(fmt.Sprintf("%s:%s", conn.GeoAddr.Host, conn.GeoAddr.Port), grpc.WithInsecure())
	if err != nil {
		log.Printf("Can't dial to grpc %s", err)
	}

	sugar.Infof("Sucseful connnect to %s:%s", conn.GeoAddr.Host, conn.GeoAddr.Port)
	authConn, err := grpc.Dial(fmt.Sprintf("%s:%s", conn.AuthAddr.Host, conn.AuthAddr.Port), grpc.WithInsecure())
	if err != nil {
		log.Printf("Can't dial to grpc %s", err)
	}
	sugar.Infof("Sucseful connnect to %s:%s", conn.AuthAddr.Host, conn.AuthAddr.Port)

	userConn, err := grpc.Dial(fmt.Sprintf("%s:%s", conn.UserAddr.Host, conn.UserAddr.Port), grpc.WithInsecure())
	if err != nil {
		log.Printf("Can't dial to grpc %s", err)
	}

	sugar.Infof("Sucseful connnect to  %s:%s", conn.UserAddr.Host, conn.UserAddr.Port)
	defer geoConn.Close()
	defer authConn.Close()
	defer userConn.Close()

	amqpConnect, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal("Can't connect to amqp ")
	}
	usersRequsts := make(chan string)
	isAllow := make(chan bool)
	go ratelimit.RateWorker(usersRequsts, isAllow)
	controller := controllers.NewControllers(sugar, userConn, authConn, geoConn)
	rp := internalMiddleware.NewReverseProxy(controller, fmt.Sprintf("http://hugo"), fmt.Sprintf(":%s", conn.HugoAddr.Port), usersRequsts, isAllow, amqpConnect)
	r.Use(rp.ReverseProxy)

	ro := router.NewRouter(r, controller)

	port := fmt.Sprintf(":%s", conn.ServConfig.Port)
	server := http.Server{
		Addr:         port,
		Handler:      ro,
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

}

//func fullGoroutineStackDump(w http.ResponseWriter, r *http.Request) {
//	// Получаем полный стек горутин
//	buf := make([]byte, 1<<20)
//	stackLen := runtime.Stack(buf, true)
//
//	// Отправляем содержимое стека в ответ на запрос
//	w.Header().Set("Content-Type", "text/plain")
//	w.WriteHeader(http.StatusOK)
//	w.Write(buf[:stackLen])
//}

//r.Route("/mycustompath", func(r chi.Router) {
//	//r.Use(jwtauth.Verifier(migrator.TokenAuth))
//	//r.Use(jwtauth.Authenticator(migrator.TokenAuth))
//	//r.Use(internalMiddleware.TokenAuthMiddleware)
//	r.Get("/pprof/profile", pprof.Profile)
//	r.Get("/pprof/trace", pprof.Trace)
//	r.Get("/pprof/", pprof.Index)
//	r.Get("/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
//	r.Get("/pprof/block", pprof.Handler("block").ServeHTTP)
//	r.Get("/pprof/cmdline", pprof.Cmdline)
//	r.Get("/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
//	r.Get("/pprof/heap", pprof.Handler("heap").ServeHTTP)
//	r.Get("/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
//	r.Get("/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
//	r.Get("/debug/pprof/goroutine", fullGoroutineStackDump)
//})
