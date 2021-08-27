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

func eventSubIndex(esub *EventSub, server *api.APIServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("EventSub route index PauseManShit\n"))
	}
}

// eventSubCallback handles POST /eventsub/callback
func eventSubCallback(esub *EventSub, server *api.APIServer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body in eventSubCallback: " + err.Error())
			return
		}
		defer r.Body.Close()

		// First of all, check if the message really came from Twitch by verifying the signature
		if !helix.VerifyEventSubNotification("", r.Header, string(body)) {
			log.Println("Received a notification, but the signature was invalid")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Read data sent in the request
		var notification eventSubNotification
		//err = json.NewDecoder(bytes.NewReader(body)).Decode(&vals)
		err = json.Unmarshal(body, &notification)
		if err != nil {
			log.Printf("Error unmarshaling incoming eventsub message: %s, request body: %s\n", err, string(body))

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// If a challenge is specified in request, respond to it
		if notification.Challenge != "" {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(notification.Challenge))
			return
		}

		esub.handleIncomingNotification(notification)
		w.WriteHeader(http.StatusOK)
	}

}

func (esub *EventSub) registerAPIRoutes(server *api.APIServer) {
	server.Router.Get("/eventsub", eventSubIndex(esub, server))
	server.Router.Post("/eventsub/callback", eventSubCallback(esub, server))
}
