package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"ikualo.com/ikualiff/internal"
	"ikualo.com/ikualiff/internal/database"
)

type VeriffResponse struct {
	Status       string
	Verification VeriffResponseVerification
}

type VeriffResponseVerification struct {
	Id           string
	Url          string
	VendorData   string
	Host         string
	Status       string
	SessionToken string
}

func RegisterFileServer() {
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/web/static/", http.StripPrefix("/web/static/", fs))
}

func RegisterLoginRoute() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		mainLayout := filepath.Join("web/templates", "layout.html")
		page := filepath.Join("web/pages", "login.html")

		tmpl, err := template.ParseFiles(mainLayout, page)
		if err != nil {
			log.Panic(err)
		}

		tmpl.ExecuteTemplate(w, "layout", nil)
	})
}

func RegisterVerifyMeRoute() {
	http.HandleFunc("/verify-me", func(w http.ResponseWriter, r *http.Request) {
		mainLayout := filepath.Join("web/templates", "auth-layout.html")
		page := filepath.Join("web/pages", "verify-me.html")

		claims, err := getRequestClaims(r)
		if err != nil {
			http.Redirect(w, r, "/login.html", http.StatusSeeOther)
			return
		}

		if err = claims.Valid(); err != nil {
			http.Redirect(w, r, "/login.html", http.StatusSeeOther)
			return
		}

		var idvUrl string
		var poaUrl string

		veriffIdvApiKey := internal.GetEnv()[internal.VERIFF_IDV_API_KEY]
		veriffIdvResponse, err := getVeriffUrl(veriffIdvApiKey)
		if err != nil {
			log.Panicf("Error getting Veriff response: %v", err)
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}

		veriffPoaApiKey := internal.GetEnv()[internal.VERIFF_POA_API_KEY]
		veriffPoaResponse, err := getVeriffUrl(veriffPoaApiKey)
		if err != nil {
			log.Panicf("Error getting Veriff response: %v", err)
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}

		err = database.CreateCustomer(
			claims.Email,
			veriffIdvResponse.Verification.Id,
			veriffIdvResponse.Verification.Url,
			veriffPoaResponse.Verification.Id,
			veriffPoaResponse.Verification.Url,
		)
		if err != nil {
			log.Panicf("User could not be created! %v", err)
			http.Redirect(w, r, "/error", http.StatusSeeOther)
			return
		}

		idvUrl = veriffIdvResponse.Verification.Url
		poaUrl = veriffPoaResponse.Verification.Url

		tmpl, err := template.ParseFiles(mainLayout, page)
		if err != nil {
			log.Fatalf("Could not parse template files: %v", tmpl)
		}

		tmpl.ExecuteTemplate(w, "auth-layout", struct {
			IdvUrl string
			PoaUrl string
		}{
			IdvUrl: idvUrl,
			PoaUrl: poaUrl,
		})
	})
}

func RegisterFinishRoute() {
	http.HandleFunc("/finish", func(w http.ResponseWriter, r *http.Request) {
		mainLayout := filepath.Join("web/templates", "layout.html")
		page := filepath.Join("web/pages", "finish.html")

		tmpl, err := template.ParseFiles(mainLayout, page)
		if err != nil {
			log.Panic(err)
		}

		tmpl.ExecuteTemplate(w, "layout", nil)
	})
}

func RegisterAnyRoute() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func getRequestClaims(r *http.Request) (*internal.Claims, error) {
	cookies := r.Cookies()

	for _, cookie := range cookies {
		if cookie.Name == "accessToken" {
			claims := internal.ParseJwtToken(cookie.Value)
			return claims, nil
		}
	}
	return nil, errors.New("accessToken not present in request.")
}

func getVeriffUrl(veriffApiKey string) (VeriffResponse, error) {
	veriffUrl := internal.GetEnv()[internal.VERIFF_URL]
	veriffCallbackUrl := internal.GetEnv()[internal.VERIFF_CALLBACK_URL]

	data, err := json.Marshal(map[string]interface{}{
		"verification": map[string]interface{}{
			"callback": veriffCallbackUrl,
		},
	})

	if err != nil {
		log.Panic(err)
	}

	request, err := http.NewRequest("POST", veriffUrl, bytes.NewBuffer(data))
	if err != nil {
		log.Panic(err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-AUTH-CLIENT", veriffApiKey)

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	var response VeriffResponse
	err = json.Unmarshal(bytes, &response)

	if err != nil {
		log.Panic(err)
		return VeriffResponse{}, err
	}
	return response, nil
}
