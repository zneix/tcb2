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

type apiServer struct {
	// baseURL ...
	baseURL string

	// bindAddress on which the HTTP server will listen on
	bindAddress string

	// listenPrefix ...
	listenPrefix string

	// router ...
	router *chi.Mux

	startTime time.Time
}

// mountRouter tries to figure out listenPrefix from server.BaseURL
func mountRouter(server *apiServer) *chi.Mux {
	if server.baseURL == "" {
		return server.router
	}

	u, err := url.Parse(server.baseURL)
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
	ur.Mount(server.listenPrefix, server.router)
	server.router = ur

	return ur
}

func (server *apiServer) Listen() {

	srv := &http.Server{
		Handler:      mountRouter(server),
		Addr:         server.bindAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("[API] Listening on %s (Prefix=%s, BaseURL=%s)\n", server.bindAddress, server.listenPrefix, server.baseURL)
	log.Fatal(srv.ListenAndServe())
}

func New(cfg config.TCBConfig) *apiServer {
	router := chi.NewRouter()

	// Strip trailing slashes from API requests
	router.Use(middleware.StripSlashes)

	server := &apiServer{
		baseURL:      cfg.BaseURL,
		bindAddress:  cfg.BindAddress,
		listenPrefix: "/",
		router:       router,
		startTime:    time.Now(),
	}

	// Handle routes
	handleMainRoutes(server)

	return server
}
