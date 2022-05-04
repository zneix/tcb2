package eventsub

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/nicklaw5/helix/v2"
	"github.com/zneix/tcb2/internal/api"
)

// routeIndex handles GET /eventsub
func (esub *EventSub) routeIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("EventSub route index PauseManShit\n"))
}

// routeCallback handles POST /eventsub/callback
func (esub *EventSub) routeCallback(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("[EventSub] Error reading request body in eventSubCallback:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// First of all, check if the message really came from Twitch by verifying the signature
	if !helix.VerifyEventSubNotification(esub.secret, r.Header, string(body)) {
		log.Println("[EventSub] Received a notification, but the signature was invalid")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Read data sent in the request
	var notification eventSubNotification
	err = json.Unmarshal(body, &notification)
	if err != nil {
		log.Printf("[EventSub] Failed to unmarshal incoming message: %s, request body: %s\n", err, string(body))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If a challenge is specified in request, respond to it
	if notification.Challenge != "" {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte(notification.Challenge))
		return
	}

	// XXX: This order of these two functions seems awfully wrong, but in stream-online-while-already-live cases it /could/ prevent panics(?)
	w.WriteHeader(http.StatusOK)
	esub.handleIncomingNotification(&notification)
}

func (esub *EventSub) registerAPIRoutes(server *api.Server) {
	server.Router.Get("/eventsub", esub.routeIndex)
	server.Router.Post("/eventsub/callback", esub.routeCallback)
}
