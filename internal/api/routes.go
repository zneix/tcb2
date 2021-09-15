package api

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/nicklaw5/helix"
	"github.com/zneix/tcb2/pkg/utils"
)

// routeIndex handles GET /api
func (server *Server) routeIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("This is the public titlechange_bot's API, most of the endpoints however are (and will be) undocumented ThreeLetterAPI TeaTime\nMore information on the GitHub repo: https://github.com/zneix/tcb2\n"))
}

// routeHealth handles GET /api/health
func (server *Server) routeHealth(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memory := fmt.Sprintf("Alloc=%v MiB, TotalAlloc=%v MiB, Sys=%v MiB, NumGC=%v",
		m.Alloc/1024/1024,
		m.TotalAlloc/1024/1024,
		m.Sys/1024/1024,
		m.NumGC)

	_, _ = w.Write([]byte(fmt.Sprintf("API Uptime: %s\nMemory: %s\n", utils.TimeSince(server.startTime), memory)))
}

// routeAuthCallback handles GET /auth/callback
func (server *Server) routeAuthCallback(helixClient *helix.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("No code provided...\n"))
			return
		}

		// Request access token for given code
		respT, err := helixClient.RequestUserAccessToken(code)
		if err != nil {
			log.Printf("Error while requesting user access token: %s\n", err)

			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal server error..."))
			return
		}

		// Requesting token was unsuccessful
		if respT.StatusCode != http.StatusOK {
			w.WriteHeader(respT.StatusCode)
			_, _ = w.Write([]byte(fmt.Sprintf("Something went wrong while requesting access token, %d %s", respT.StatusCode, respT.ErrorMessage)))
			return
		}

		// Validate obtained token for extra data (no way it can be invalid)
		_, respV, _ := helixClient.ValidateToken(respT.Data.AccessToken)

		log.Printf("%# v\n", respV)

		// All went good
		_, _ = w.Write([]byte("KKona\n"))
	}
}

func registerMainRoutes(server *Server, helixClient *helix.Client) {
	server.Router.Get("/api", server.routeIndex)
	server.Router.Get("/api/health", server.routeHealth)
	server.Router.Get(authCallbackRoute, server.routeAuthCallback(helixClient))
}
