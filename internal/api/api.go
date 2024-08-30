package api

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"ikualo.com/ikualiff/internal"
)

func RegisterLoginRoute() {
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		log.Print("New login request!")

		r.ParseForm()
		email := r.Form.Get("email")
		code := r.Form.Get("code")

		apiKey := internal.GetEnv()[internal.API_KEY]

		if code == apiKey {
			cookie := generateJwtTokenCookie(email)

			http.SetCookie(w, cookie)

			w.Header().Add("HX-Redirect", "/verify-me")
			w.WriteHeader(http.StatusOK)
			return
		}

		errorComponent := filepath.Join("web/components", "errorCard.html")
		tmpl, err := template.ParseFiles(errorComponent)
		if err != nil {
			log.Fatalf("Error loading error component  template: %v", err)
		}

		tmpl.Execute(w, "Oh... ese código no es el correcto.")
	})
}

func generateJwtTokenCookie(email string) (*http.Cookie) {
	token := internal.GenerateJwt(email)

	cookie := http.Cookie{}
	cookie.Name = "accessToken"
	cookie.Value = token
	cookie.Expires = time.Now().Add(10 * time.Minute)
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Secure = true
	cookie.HttpOnly = true
	cookie.Path = "/"

	return &cookie
}
