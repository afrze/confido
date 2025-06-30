package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Initialize Viper
	initConfig()

	env := viper.GetString("app.env")
	port := viper.GetString("server.port")

	handler := buildRouter(env, port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		log.Printf("Confido is running on http://localhost:%s (env=%s)", port, env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)

}

func initConfig() {
	// Defaults
	viper.SetDefault("app.env", "development")
	viper.SetDefault("server.port", 8040)

	// Pull from real environment variables
	viper.AutomaticEnv()

	// Load .config.yaml if present
	viper.SetConfigName(".config")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
}

func buildRouter(env, port string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "pong")
	})

	var root http.Handler

	if env == "production" {
		root = http.FileServer(http.Dir("./client/dist"))
	} else {
		root = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(w, r)
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "Confido is running in %s mode on port %s\n", env, port)
		})
	}

	mux.Handle("/", root)
	return mux
}
