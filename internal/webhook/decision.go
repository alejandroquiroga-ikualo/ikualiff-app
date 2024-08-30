package webhook

import (
	"io"
	"log"
	"net/http"

	"ikualo.com/ikualiff/internal/database"
)

func CreateDecisionTable() {
	_, err := database.CreateTable(`
		CREATE TABLE IF NOT EXISTS decision (
			id SERIAL PRIMARY KEY,
			decision JSON NOT NULL
		)
	`)

	if err != nil {
		log.Fatalf("Couldn't create decision table! %v", err)
	}
}

func RegisterDecisionRoute() {
	http.HandleFunc("/webhook/decisions", func(w http.ResponseWriter, r *http.Request) {
		log.Print("New decision!")

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("Error reading request body!: %v", err)
		}

		json := string(bytes)
		err = saveDecision(json)
		if err != nil {
			log.Panicf("Error saving decision! %v", err)
		}
	})
}

func saveDecision(json string) error {
	log.Print("Saving new decision...")
	return database.Exec(`
		INSERT INTO decision (decision)
		VALUES ($1);
	`, json)
}
