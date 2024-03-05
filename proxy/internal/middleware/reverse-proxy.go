package middleware

import (
	"net/http"
	"net/http/httputil"
	"net/url"
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
}

func NewReverseProxy(controller *controllers.Controllers, host, port string) *ReverseProxy {
	//queue, err := rabbitmq.NewRabbitMQ(conn)
	//if err != nil {
	//	log.Fatalln("Can't create ")
	//}
	return &ReverseProxy{
		controller: controller,
		host:       host,
		port:       port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/auth") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/user") || strings.HasPrefix(r.URL.Path, "/geo") || strings.HasPrefix(r.URL.Path, "/swagger") {
			token := r.Header.Get("Authorziation")
			isAllow := rp.controller.AuthController.Verif(token)
			if !isAllow {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			//if strings.HasPrefix(r.URL.Path, "/api/geo") {
			//	token := r.Header.Get("Authorization")
			//	splitedToken := strings.Split(token, " ")
			//	rp.req <- splitedToken[1]
			//	select {
			//	case allow := <-rp.allow:
			//		if !allow {
			//			err := rp.rabbitmq.Publish("to_many", []byte(splitedToken[1]))
			//			if err != nil {
			//				http.Error(w, "Internal error", http.StatusInternalServerError)
			//				return
			//			}
			//			http.Error(w, "To many requests", http.StatusTooManyRequests)
			//			return
			//		}
			//	}
			//
			//}
			//if isAuth := rp.controller.AuthController.Verif(r.Header.Get("Authorization")); !isAuth {
			//	w.WriteHeader(http.StatusUnauthorized)
			//	return
			//}
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
