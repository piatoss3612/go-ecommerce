package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// set up cors
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Post("/api/payment-intent", app.GetPaymentIntent)

	mux.Get("/api/widget/{id}", app.GetWidgetById)

	mux.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)

	mux.Post("/api/authenticate", app.CreateAuthToken)

	mux.Post("/api/is-authenticated", app.CheckAuthentication)

	// create new mux and apply middleware to it
	// routes starting with /api/admin will be grouped together and protected by middleware
	mux.Route("/api/admin", func(mux chi.Router) {
		mux.Use(app.Auth)

		// should be authenticated to post virtual terminal request
		mux.Post("/virtual-terminal-succeeded", app.VirtualTerminalPaymentSucceeded)

		mux.Post("/all-sales", app.AllSales)
		mux.Post("/get-sales/{id}", app.GetSale)

		mux.Post("/all-subscriptions", app.AllSubscriptions)

		mux.Post("/refund", app.RefundCharge)
	})

	mux.Post("/api/forgot-password", app.SendPasswordResetEmail)

	mux.Post("/api/reset-password", app.ResetPassword)

	return mux
}
