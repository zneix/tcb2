package api

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/zneix/tcb2/pkg/utils"
)

// index handles GET /
func index(server *APIServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("This is the public titlechange_bot's API, most of the endpoints however are (and will be) undocumented ThreeLetterAPI TeaTime\nMore information on the GitHub repo: https://github.com/zneix/tcb2\n"))
	}
}

// health handles GET /health
func health(server *APIServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		memory := fmt.Sprintf("Alloc=%v MiB, TotalAlloc=%v MiB, Sys=%v MiB, NumGC=%v",
			m.Alloc/1024/1024,
			m.TotalAlloc/1024/1024,
			m.Sys/1024/1024,
			m.NumGC)

		w.Write([]byte(fmt.Sprintf("API Uptime: %s\nMemory: %s\n", utils.TimeSince(server.startTime), memory)))
	}
}

func registerMainRoutes(server *APIServer) {
	server.Router.Get("/", index(server))
	server.Router.Get("/health", health(server))
}
