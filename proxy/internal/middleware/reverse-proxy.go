package middleware

import (
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"proxy/internal/infra/queue"
	"proxy/internal/infra/rabbitmq"
	controllers "proxy/internal/modules"
	"strings"
)

// rabbitmq   queue.MessageQueuer
type ReverseProxy struct {
	controller *controllers.Controllers
	host       string
	port       string
	req        chan string
	allow      chan bool
	rabbitmq   queue.MessageQueuer
}

func NewReverseProxy(controller *controllers.Controllers, host, port string, req chan string, allow chan bool, conn *amqp.Connection) *ReverseProxy {
	queue, err := rabbitmq.NewRabbitMQ(conn)
	if err != nil {
		log.Fatalln("Can't create ")
	}
	return &ReverseProxy{
		controller: controller,
		host:       host,
		port:       port,
		req:        req,
		allow:      allow,
		rabbitmq:   queue,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/auth") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/user") || strings.HasPrefix(r.URL.Path, "/api") || strings.HasPrefix(r.URL.Path, "/swagger") {
			token := r.Header.Get("Authorization")
			isAllow := rp.controller.AuthController.Verif(token)
			if !isAllow {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if strings.HasPrefix(r.URL.Path, "/api") {
				jwt := r.Header.Get("Authorization")
				rp.req <- jwt
				select {
				case allow := <-rp.allow:
					if !allow {
						err := rp.rabbitmq.Publish("GeoRateLimit", []byte(jwt))
						if err != nil {
							log.Println(err)
							http.Error(w, "Internal error", http.StatusInternalServerError)
							return
						}
						http.Error(w, "To many requests", http.StatusTooManyRequests)
						return
					}
				}

			}
			next.ServeHTTP(w, r)
			return
		}
		targetURL, err := url.Parse(rp.host + rp.port)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		proxy.ServeHTTP(w, r)
	})
}
