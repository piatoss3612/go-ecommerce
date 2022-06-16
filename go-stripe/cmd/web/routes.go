package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad) // middleware

	// home page
	mux.Get("/", app.Home)

	// virtual terminal page protected by middleware
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.Auth)
		mux.Get("/virtual-terminal", app.VirtualTerminal)

		mux.Get("/all-sales", app.AllSales)
		mux.Get("/sales/{id}", app.ShowSale)

		mux.Get("/all-subscriptions", app.AllSubscriptions)
		mux.Get("/subscriptions/{id}", app.ShowSubscription)

		mux.Get("/all-users", app.AllUsers)
		mux.Get("/all-users/{id}", app.OneUser)
	})

	// widget page
	mux.Get("/widget/{id}", app.ChargeOnce)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/receipt", app.Receipt)

	// subscription page
	mux.Get("/plans/bronze", app.BronzePlan)
	mux.Get("/receipt/bronze", app.BronzePlanReceipt)

	// authentication page
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)

	mux.Get("/logout", app.Logout)

	mux.Get("/forgot-password", app.ForgotPassword)
	mux.Get("/reset-password", app.ShowResetPassword)

	fileServer := http.FileServer(http.Dir("./static")) // use file system
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
