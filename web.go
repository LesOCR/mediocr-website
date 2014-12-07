package main

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gopkg.in/unrolled/render.v1"
)

/* Request handlers */

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r, "home")
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, r, mux.Vars(r)["page"])
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	redirectError(w, r, http.StatusNotFound)
}

/* Error redirection functions */

// redirectErrorDesc shows a flash message additionally to an error
func redirectErrorDesc(w http.ResponseWriter, r *http.Request, errorCode int,
	description string) {
	getSession(r).AddFlash(description, "error_description")
	redirectError(w, r, errorCode)
}

// redirectError redirects the user to an error page
func redirectError(w http.ResponseWriter, r *http.Request, errorCode int) {
	getSession(r).AddFlash(errorCode, "error_code")
	saveSession(w, r)
	http.Redirect(w, r, "/error", http.StatusSeeOther)
}

/* Context functions */

func contextMiddleware(rn *render.Render, store *sessions.CookieStore) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		context.Set(r, "render", rn)
		context.Set(r, "store", store)
		next(w, r)
	}
}

// getRender flushes the session to the browser and returns the renderer
func getRender(w http.ResponseWriter, r *http.Request) *render.Render {
	return context.Get(r, "render").(*render.Render)
}

func getSession(r *http.Request) *sessions.Session {
	store := context.Get(r, "store").(*sessions.CookieStore)
	// We ignore the error because store.Get() always returns a session
	session, _ := store.Get(r, "mediocr")
	return session
}

func saveSession(w http.ResponseWriter, r *http.Request) error {
	return sessions.Save(r, w)
}
