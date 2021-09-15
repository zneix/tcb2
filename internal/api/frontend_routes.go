package api

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/zneix/tcb2/internal/config"
)

var scopes = []string{}

type modelFrontendIndex struct {
	ListenPrefix string
	LoginStatus  string
	LoginURI     string
}

// routeFrontendIndex handles GET /
func (server *Server) routeFrontendIndex(cfg *config.TCBConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		indexTemplate, err := template.ParseFiles("web/index.html")
		if err != nil {
			log.Fatalf("[API-Frontend] Error parsing index template file: %s\n", err)
		}

		// Execute html template and write it as the response
		data := &modelFrontendIndex{
			ListenPrefix: server.listenPrefix,
			LoginStatus:  "TODO: Add login handler",
			LoginURI:     server.twitchLoginURI,
		}
		indexTemplate.Execute(w, data)
	}
}

// routeStatic handles GET /static/*
func (server *Server) routeStatic(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "web/src"))
	fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
	fs.ServeHTTP(w, r)
}

func registerFrontendRoutes(server *Server, cfg *config.TCBConfig) {
	// Static content
	server.Router.Get("/static/*", server.routeStatic)

	// Routes using templates
	server.Router.Get("/", server.routeFrontendIndex(cfg))
}
