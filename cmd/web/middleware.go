package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/surakshith-suvarna/bookings/internal/helpers"
)

//NoSurf adds CSRF token for every POST request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

//Auth verifies if a user is authenticated
func Auth(next http.Handler) http.Handler {
	//Returns a handler by calling an annonynous function.
	//This is possible because we are returning something which has access to responsewriter and request
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			app.Session.Put(r.Context(), "error", "Please login first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		//Once completed, Go to next http request
		next.ServeHTTP(w, r)

	})
}
