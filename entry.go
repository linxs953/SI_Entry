package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"

	"entry/internal/config"
	"entry/internal/handler"
	"entry/internal/svc"
)

var configFile = flag.String("f", "etc/entry-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(corsMiddleware)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Errorf("[CORS] Handling request: %s %s (Origin: %s)", r.Method, r.URL.Path, r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next(w, r)
	}
	// return func(w http.ResponseWriter, r *http.Request) {
	// 	origin := r.Header.Get("Origin")
	// 	if origin != "" {
	// 		w.Header().Set("Access-Control-Allow-Origin", origin)
	// 		w.Header().Set("Vary", "Origin")
	// 	}

	// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
	// 	w.Header().Set("Access-Control-Allow-Headers", strings.Join([]string{
	// 		"Accept",
	// 		"Authorization",
	// 		"Content-Type",
	// 		"X-CSRF-Token",
	// 		"X-Requested-With",
	// 	}, ", "))

	// 	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// 	w.Header().Set("Access-Control-Max-Age", "86400")

	// 	if r.Method == "OPTIONS" {
	// 		w.WriteHeader(http.StatusNoContent)
	// 		return
	// 	}

	// 	next(w, r)
	// }
}
