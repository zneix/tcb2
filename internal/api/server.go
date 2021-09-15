package api

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nicklaw5/helix"
	"github.com/zneix/tcb2/internal/config"
)

const (
	authCallbackRoute = "/auth/callback"
)

type Server struct {
	// Router ...
	Router *chi.Mux

	// BaseURL ...
	BaseURL string

	// bindAddress on which the HTTP server will listen on
	bindAddress string

	// listenPrefix ...
	listenPrefix string

	// twitchLoginURI ...
	twitchLoginURI string

	startTime time.Time
}

// mountRouter tries to figure out listenPrefix from server.BaseURL
func mountRouter(server *Server) *chi.Mux {
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

// Listen starts to listen on configured bindAddress (blocking)
func (server *Server) Listen() {
	srv := &http.Server{
		Handler:      mountRouter(server),
		Addr:         server.bindAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("[API] Listening on %s (Prefix=%s, BaseURL=%s)\n", server.bindAddress, server.listenPrefix, server.BaseURL)
	log.Fatal(srv.ListenAndServe())
}

func New(cfg *config.TCBConfig, helixClient *helix.Client) *Server {
	router := chi.NewRouter()

	// Strip trailing slashes from API requests
	router.Use(middleware.StripSlashes)

	// Apply the RedirectURI to helixClient
	redirectURI := cfg.BaseURL + authCallbackRoute
	helixClient.SetRedirectURI(redirectURI)

	// Figure out twitchLoginURI
	const responseType = "code"
	const forceVerify = "true"
	encodedScopes := url.PathEscape(strings.Join(scopes, " "))

	twitchLoginURL, _ := url.Parse("https://id.twitch.tv/oauth2/authorize")
	twitchLoginURLVariables := &url.Values{}
	twitchLoginURLVariables.Set("client_id", cfg.TwitchClientID)
	twitchLoginURLVariables.Set("redirect_uri", redirectURI)
	twitchLoginURLVariables.Set("response_type", responseType)
	twitchLoginURLVariables.Set("scope", encodedScopes)
	twitchLoginURLVariables.Set("force_verify", forceVerify)
	twitchLoginURL.RawQuery = twitchLoginURLVariables.Encode()

	// Create new Server instance
	server := &Server{
		Router:         router,
		BaseURL:        cfg.BaseURL,
		bindAddress:    cfg.BindAddress,
		listenPrefix:   "/",
		startTime:      time.Now(),
		twitchLoginURI: twitchLoginURL.String(),
	}

	// Handle routes
	registerMainRoutes(server, helixClient)
	registerFrontendRoutes(server, cfg)

	return server
}
