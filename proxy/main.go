package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"hugoproxy-main/proxy/controllers"
	internalMiddleware "hugoproxy-main/proxy/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	time.Sleep(time.Second * 15)
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	r := chi.NewRouter()
	r.Use(middleware.DefaultLogger)
	client := http.Client{}
	geoConn, err := grpc.Dial("geo:1234", grpc.WithInsecure())
	if err != nil {
		log.Printf("Can't dial to grpc %s", err)
	}
	authConn, err := grpc.Dial("auth:1235", grpc.WithInsecure())
	if err != nil {
		log.Printf("Can't dial to grpc %s", err)
	}
	userConn, err := grpc.Dial("user_service:1237", grpc.WithInsecure())
	if err != nil {
		log.Printf("Can't dial to grpc %s", err)
	}
	sugar.Infof("Sucseful connnect to geo:1234")
	sugar.Infof("Sucseful connnect to auth:1235")
	sugar.Infof("Sucseful connnect to user_service:1237")
	defer geoConn.Close()
	defer authConn.Close()
	defer userConn.Close()
	controller := controllers.NewControllers(sugar, client, geoConn, authConn, userConn)
	rp := internalMiddleware.NewReverseProxy(controller, "http://hugo", ":1313")
	r.Use(rp.ReverseProxy)
	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", controller.AuthController.Register)
			r.Post("/login", controller.AuthController.Login)
		})
		r.Route("/user_service", func(r chi.Router) {
			r.Post("/profile", controller.UserController.Profile)
			r.Post("/list", controller.UserController.List)
		})
		r.Route("/address", func(r chi.Router) {
			r.Post("/search", controller.SearchController.GetSearch)
			r.Post("/geocode", controller.SearchController.GetGeo)
		})

	})
	r.Route("/swagger", func(r chi.Router) {
		r.Get("/index", controller.SwagController.GetSwaggerHtml)
		r.Get("/swagger", controller.SwagController.GetSwaggerJson)
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

}

//log.Println("Server exiting")
//router := chi.NewRouter()
//router.Use()
//func Auth(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//	})
//}

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
//	//r.Use(jwtauth.Verifier(storage.TokenAuth))
//	//r.Use(jwtauth.Authenticator(storage.TokenAuth))
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
//{"id":123,"name":"123","phone":"123","email":"123","password":"123"}
//if configrpc == "rpc" {
//	r.Route("/api", func(r chi.Router) {
//		//r.Use(jwtauth.Verifier(storage.TokenAuth))
//		//r.Use(jwtauth.Authenticator(storage.TokenAuth))
//		//r.Use(internalMiddleware.TokenAuthMiddleware)
//		r.Route("/address", func(r chi.Router) {
//
//			r.Post("/search", controller.SearchController.GetSearch)
//			r.Post("/geocode", controller.SearchController.GetGeoCode)
//		})
//	})
//} else if configrpc == "json-rpc" {
//	r.Route("/api", func(r chi.Router) {
//		//r.Use(jwtauth.Verifier(storage.TokenAuth))
//		//r.Use(jwtauth.Authenticator(storage.TokenAuth))
//		//r.Use(internalMiddleware.TokenAuthMiddleware)
//		r.Route("/address", func(r chi.Router) {
//
//			r.Post("/search", controller.SearchControllerJsonRpc.GetSearch)
//			r.Post("/geocode", controller.SearchControllerJsonRpc.GetGeoCode)
//		})
//	})
