package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad) // middleware

	mux.Get("/", app.Home)
	mux.Get("/virtual-terminal", app.VirtualTerminal)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/receipt", app.Receipt)
	mux.Get("/widget/{id}", app.ChargeOnce)

	fileServer := http.FileServer(http.Dir("./static")) // use file system
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
