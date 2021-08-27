package api

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zneix/tcb2/internal/config"
)

type APIServer struct {
	// Router ...
	Router *chi.Mux

	// BaseURL ...
	BaseURL string

	// bindAddress on which the HTTP server will listen on
	bindAddress string

	// listenPrefix ...
	listenPrefix string

	startTime time.Time
}

// mountRouter tries to figure out listenPrefix from server.BaseURL
func mountRouter(server *APIServer) *chi.Mux {
	if server.BaseURL == "" {
		return server.Router
	}

	u, err := url.Parse(server.BaseURL)
	if err != nil {
		log.Fatalln("[API] Error mounting router: " + err.Error())
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		log.Fatalln("[API] Scheme must be included in Base URL")
	}

	if u.Path != "" {
		server.listenPrefix = u.Path
	}
	ur := chi.NewRouter()
	ur.Mount(server.listenPrefix, server.Router)
	server.Router = ur

	return ur
}

func (server *APIServer) Listen() {

	srv := &http.Server{
		Handler:      mountRouter(server),
		Addr:         server.bindAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("[API] Listening on %s (Prefix=%s, BaseURL=%s)\n", server.bindAddress, server.listenPrefix, server.BaseURL)
	log.Fatal(srv.ListenAndServe())
}

func New(cfg config.TCBConfig) *APIServer {
	router := chi.NewRouter()

	// Strip trailing slashes from API requests
	router.Use(middleware.StripSlashes)

	server := &APIServer{
		Router:       router,
		BaseURL:      cfg.BaseURL,
		bindAddress:  cfg.BindAddress,
		listenPrefix: "/",
		startTime:    time.Now(),
	}

	// Handle routes
	registerMainRoutes(server)

	return server
}
