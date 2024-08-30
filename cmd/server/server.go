package main

import (
	"log"
	"net/http"

	"ikualo.com/ikualiff/internal"
	"ikualo.com/ikualiff/internal/api"
	"ikualo.com/ikualiff/internal/database"
	"ikualo.com/ikualiff/internal/web"
	"ikualo.com/ikualiff/internal/webhook"
)

func main() {
	internal.GetEnv()

	web.RegisterFileServer()

	api.RegisterLoginRoute()

	web.RegisterLoginRoute()
	web.RegisterVerifyMeRoute()
	web.RegisterAnyRoute()

	webhook.CreateDecisionTable()
	webhook.RegisterDecisionRoute()

	webhook.CreateEventTable()
	webhook.RegisterEventRoute()

	database.CreateUserTable()

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
