package api

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/zneix/tcb2/pkg/utils"
)

// routeIndex handles GET /
func (server *Server) routeIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("This is the public titlechange_bot's API, most of the endpoints however are (and will be) undocumented ThreeLetterAPI TeaTime\nMore information on the GitHub repo: https://github.com/zneix/tcb2\n"))
}

// routeHealth handles GET /routeHealth
func (server *Server) routeHealth(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memory := fmt.Sprintf("Alloc=%v MiB, TotalAlloc=%v MiB, Sys=%v MiB, NumGC=%v",
		m.Alloc/1024/1024,
		m.TotalAlloc/1024/1024,
		m.Sys/1024/1024,
		m.NumGC)

	// _, _ = w.Write([]byte(fmt.Sprintf("API Uptime: %s\nMemory: %s\n", utils.TimeSince(server.startTime), memory)))
	fmt.Fprintf(w, "API Uptime: %s\nMemory: %s\n", utils.TimeSince(server.startTime), memory)
}

func registerMainRoutes(server *Server) {
	server.Router.Get("/", server.routeIndex)
	server.Router.Get("/health", server.routeHealth)
}
