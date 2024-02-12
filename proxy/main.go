package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"hugoproxy-main/proxy/controllers"
	internalMiddleware "hugoproxy-main/proxy/middleware"
	"log"
	"net/http"
	"net/http/pprof"
	rpc "net/rpc"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can't load config from env")
	}
	time.Sleep(time.Second * 15)
	configrpc := os.Getenv("RPC_PROVIDER")
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	client := http.Client{}
	rpc, err := rpc.DialHTTP("tcp", "serv:1234")
	if err != nil {
		log.Println("Can't connect to rpc")
	}
	grpc, err := grpc.Dial("serv:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Can't connect to grpc")
	}
	defer grpc.Close()
	controller := controllers.NewControllers(sugar, client, rpc, grpc)
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
	if configrpc == "rpc" {
		r.Route("/api", func(r chi.Router) {
			//r.Use(jwtauth.Verifier(storage.TokenAuth))
			//r.Use(jwtauth.Authenticator(storage.TokenAuth))
			//r.Use(internalMiddleware.TokenAuthMiddleware)
			r.Route("/address", func(r chi.Router) {

				r.Post("/search", controller.SearchController.GetSearch)
				r.Post("/geocode", controller.SearchController.GetGeoCode)
			})
		})
	} else if configrpc == "json-rpc" {
		r.Route("/api", func(r chi.Router) {
			//r.Use(jwtauth.Verifier(storage.TokenAuth))
			//r.Use(jwtauth.Authenticator(storage.TokenAuth))
			//r.Use(internalMiddleware.TokenAuthMiddleware)
			r.Route("/address", func(r chi.Router) {

				r.Post("/search", controller.SearchControllerJsonRpc.GetSearch)
				r.Post("/geocode", controller.SearchControllerJsonRpc.GetGeoCode)
			})
		})
	} else if configrpc == "grpc" {
		r.Route("/api", func(r chi.Router) {
			//r.Use(jwtauth.Verifier(storage.TokenAuth))
			//r.Use(jwtauth.Authenticator(storage.TokenAuth))
			//r.Use(internalMiddleware.TokenAuthMiddleware)
			r.Route("/address", func(r chi.Router) {

				r.Post("/search", controller.SearchControllergRpc.GetSearch)
				r.Post("/geocode", controller.SearchControllergRpc.GetGeo)
			})
		})
	}

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
