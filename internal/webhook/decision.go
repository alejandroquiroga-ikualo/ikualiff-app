package webhook

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
	"ikualo.com/ikualiff/internal/database"
)

type Decision struct {
	Status string
	Verification Verification
}

type Verification struct {
	Id string
}

func CreateDecisionTable() {
	_, err := database.CreateTable(`
		CREATE TABLE IF NOT EXISTS decision (
			id SERIAL PRIMARY KEY,
			veriffSessionId UUID NOT NULL,
			decision JSON NOT NULL
		)
	`)

	if err != nil {
		log.Fatalf("Couldn't create decision table! %v", err)
	}
}

func RegisterDecisionRoute(app *newrelic.Application) {
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/webhook/decisions", func(w http.ResponseWriter, r *http.Request) {
		log.Print("New decision!")

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("Error reading request body!: %v", err)
		}

		j := string(bytes)

		var data Decision
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			log.Fatalf("Error unmarshaling decision data! %v", err)
		}

		veriffSessionId := data.Verification.Id

		err = saveDecision(veriffSessionId, j)
		if err != nil {
			log.Panicf("Error saving decision! %v", err)
		}
	}))
}

func saveDecision(veriffSessionId string, json string) error {
	log.Print("Saving new decision...")
	return database.Exec(`
		INSERT INTO decision (veriffSessionId, decision)
		VALUES ($1, $2);
	`, veriffSessionId, json)
}
