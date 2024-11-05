package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gopkg.in/yaml.v2"

	"oauth-demo/internal/config"
	"oauth-demo/internal/handlers"
	"oauth-demo/internal/middleware"
)

var (
	cfg   config.Config
	store = sessions.NewCookieStore([]byte("secret-key-replace-this"))
)

func main() {
	// Load configuration
	configFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handlers and middleware with dependencies
	handlers.InitHandlers(store, &cfg)
	middleware.InitMiddleware(store)

	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/", handlers.HandleHome).Methods("GET")
	r.HandleFunc("/login", handlers.HandleLogin).Methods("GET")
	r.HandleFunc("/callback", handlers.HandleCallback).Methods("GET")

	// Protected routes
	protected := r.PathPrefix("/protected").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("", handlers.HandleProtected).Methods("GET")

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, r))
}
