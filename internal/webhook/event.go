package webhook

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
	"ikualo.com/ikualiff/internal/database"
)

type Event struct {
	Id string
	AttemptId       string
	Feature         string
	Code            int
	Action          string
	VendorData      string
	EndUserId       string
}

func CreateEventTable() {
	_, err := database.CreateTable(`
		CREATE TABLE IF NOT EXISTS event (
			id SERIAL PRIMARY KEY,
			veriffSessionId UUID NOT NULL,
			attemptId UUID NOT NULL,
			feature TEXT NOT NULL,
			code SMALLINT NOT NULL,
			action TEXT NOT NULL,
			vendorData TEXT,
			endUserId TEXT
		);
	`)

	if err != nil {
		log.Fatalf("Couldn't create event table! %v", err)
	}
}

func RegisterEventRoute(app *newrelic.Application) {
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/webhook/events", func(w http.ResponseWriter, r *http.Request) {
		log.Print("New event!")

		var event Event
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			log.Fatalf("Error decoding event: %v", err)
		}

		err = saveEvent(&event)
		if err != nil {
			log.Panicf("Error saving event! %v", err)
		}
	}))
}

func saveEvent(event *Event) error {
	log.Print("Saving new event...")
	return database.Exec(`
		INSERT INTO event (veriffSessionId, attemptId, feature, code, action, vendorData, endUserId)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`,
		event.Id,
		event.AttemptId,
		event.Feature,
		event.Code,
		event.Action,
		event.VendorData,
		event.EndUserId,
	)
}
