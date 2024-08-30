package main

import (
	"log"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
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

	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("ikualiff"),
		newrelic.ConfigLicense("5022213ea7cd500d4ee51f1daa4e4aeeFFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	webhook.CreateDecisionTable()
	webhook.RegisterDecisionRoute(app)

	webhook.CreateEventTable()
	webhook.RegisterEventRoute(app)

	database.CreateCustomerTable()

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
