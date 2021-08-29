package eventsub

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nicklaw5/helix"
	"github.com/zneix/tcb2/internal/api"
)

// eventSubIndex handles GET /eventsub
func eventSubIndex(esub *EventSub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("EventSub route index PauseManShit\n"))
	}
}

// eventSubCallback handles POST /eventsub/callback
func eventSubCallback(esub *EventSub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("[EventSub] Error reading request body in eventSubCallback: " + err.Error())
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

		esub.handleIncomingNotification(&notification)
		w.WriteHeader(http.StatusOK)
	}
}

func (esub *EventSub) registerAPIRoutes(server *api.Server) {
	server.Router.Get("/eventsub", eventSubIndex(esub))
	server.Router.Post("/eventsub/callback", eventSubCallback(esub))
}
