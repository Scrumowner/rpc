package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/" && r.Method == http.MethodGet {
			fmt.Fprintf(w, "Hello from api")
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/debug/pprof/") {
			url, _ := url.Parse(r.Host + r.URL.Path)
			debugProxy := httputil.NewSingleHostReverseProxy(url)
			debugProxy.ServeHTTP(w, r)
		}
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/swagger") {
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
