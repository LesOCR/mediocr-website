package main

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gopkg.in/unrolled/render.v1"
	"gopkg.in/unrolled/secure.v1"
)

const (
	contentsPath = "./contents"
	devMode      = false

	ocrMaxFileSize = 512 * 1024 // 512 KiB
)

func main() {
	// Allow users only from mediocr.io, and a few more security tricks
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:       []string{"mediocr.io"},
		STSSeconds:         315360000,
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		IsDevelopment:      devMode,
	})
	// Session cookie store
	store := sessions.NewCookieStore(
		[]byte(os.Getenv("SESSION_AUTHENTICATION_KEY")),
		[]byte(os.Getenv("SESSION_ENCRYPTION_KEY")),
	)
	store.Options = &sessions.Options{
	//Path: "/",
	//Secure:   true,
	//HttpOnly: true,
	}

	// Gorilla router, render package, negroni middlewares (including the
	// magical security middleware and the session store)
	r := mux.NewRouter()
	rn := render.New(render.Options{IsDevelopment: devMode})
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.Use(contextMiddleware(rn, store))

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/try", uploadHandler).Methods("POST")
	r.HandleFunc("/{page:[\\w]+}", pageHandler)
	r.HandleFunc("/{_:.*}", notFoundHandler)

	n.UseHandler(r)
	n.Run(":" + os.Getenv("PORT"))
}
