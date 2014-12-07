package main

import (
	"fmt"
	"net/http"

	"github.com/GoIncremental/negroni-sessions"
	"github.com/gorilla/mux"
	"gopkg.in/unrolled/render.v1"
)

/* Request handlers */

func homeHandler(rn *render.Render) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Host)
		renderPage(w, r, "home", rn)
	}
}

func pageHandler(rn *render.Render) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		renderPage(w, r, mux.Vars(r)["page"], rn)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	redirectError(w, r, http.StatusNotFound)
}

/* Error redirection functions */

// redirectErrorDesc shows a flash message additionally to an error
func redirectErrorDesc(w http.ResponseWriter, r *http.Request, errorCode int,
	description string) {
	sessions.GetSession(r).AddFlash(description, "error_description")
	redirectError(w, r, errorCode)
}

// redirectError redirects the user to an error page
func redirectError(w http.ResponseWriter, r *http.Request, errorCode int) {
	sessions.GetSession(r).AddFlash(errorCode, "error_code")
	http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
}
