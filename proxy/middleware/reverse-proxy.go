package middleware

import (
	"hugoproxy-main/proxy/controllers"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ReverseProxy struct {
	controller *controllers.Controllers
	host       string
	port       string
}

func NewReverseProxy(controller *controllers.Controllers, host, port string) *ReverseProxy {
	return &ReverseProxy{
		controller: controller,
		host:       host,
		port:       port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/auth/") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/debug/pprof/") {
			url, _ := url.Parse(r.Host + r.URL.Path)
			debugProxy := httputil.NewSingleHostReverseProxy(url)
			debugProxy.ServeHTTP(w, r)
		}
		if strings.HasPrefix(r.URL.Path, "/api/user_service") || strings.HasPrefix(r.URL.Path, "/api/geo") || strings.HasPrefix(r.URL.Path, "/swagger") {
			if isAuth := rp.controller.AuthController.Verif(r.Header.Get("Authorization")); !isAuth {
				w.WriteHeader(http.StatusUnauthorized)
				return
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
