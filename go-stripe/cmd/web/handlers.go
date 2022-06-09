package main

import (
	"fmt"
	"myapp/internal/cards"
	"myapp/internal/models"
	"myapp/internal/urlsigner"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// display home page
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", nil, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

// display virtual termial page
func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "terminal", nil); err != nil {
		app.errorLog.Println(err)
	}
}

// handle virtual terminal transaction data and redirect to recipt page
func (app *application) VirtualTerminalPaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	txnData, err := app.GetTransactionData(r)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// create a new transaction
	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
		ExpiryMonth:         int(txnData.ExpiryMonth),
		ExpiryYear:          int(txnData.ExpiryYear),
		BankReturnCode:      txnData.BankReturnCode,
		TransactionStatusID: 2,
	}

	_, err = app.SaveTransaction(txn)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// wirte transaction data to session,
	// and then redirect user to new page to prevent from resubmit form
	app.Session.Put(r.Context(), "receipt", txnData)                      // put data to session
	http.Redirect(w, r, "/virtual-terminal-receipt", http.StatusSeeOther) // redirect
}

// display virtual terminal receipt page
func (app *application) VirtualTerminalReceipt(w http.ResponseWriter, r *http.Request) {
	txn := app.Session.Get(r.Context(), "receipt").(TransactionData) // grap data from session
	data := make(map[string]any)
	data["txn"] = txn
	app.Session.Remove(r.Context(), "receipt")                                    // remove data from session
	if err := app.renderTemplate(w, r, "virtual-terminal-receipt", &templateData{ // render receipt page
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

// display the page to buy on widget
func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	widgetID, err := strconv.Atoi(id)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(map[string]any)
	data["widget"] = widget

	if err := app.renderTemplate(w, r, "buy-once", &templateData{
		Data: data,
	}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

// handle widget transaction data and redirect to recipt page
func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	widgetID, err := strconv.Atoi(r.Form.Get("product_id"))
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	txnData, err := app.GetTransactionData(r)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// create a new customer
	customerID, err := app.SaveCustomer(txnData.FirstName, txnData.LastName, txnData.Email)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// create a new transaction
	txn := models.Transaction{
		Amount:              txnData.PaymentAmount,
		Currency:            txnData.PaymentCurrency,
		LastFour:            txnData.LastFour,
		PaymentIntent:       txnData.PaymentIntentID,
		PaymentMethod:       txnData.PaymentMethodID,
		ExpiryMonth:         int(txnData.ExpiryMonth),
		ExpiryYear:          int(txnData.ExpiryYear),
		BankReturnCode:      txnData.BankReturnCode,
		TransactionStatusID: 2,
	}

	txnID, err := app.SaveTransaction(txn)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// create a new order
	order := models.Order{
		WidgetID:      widgetID,
		TransactionID: txnID,
		CustomerID:    customerID,
		StatusID:      1,
		Quantity:      1,
		Amount:        txnData.PaymentAmount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = app.SaveOrder(order)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// wirte transaction data to session,
	// and then redirect user to new page to prevent from resubmit form
	app.Session.Put(r.Context(), "receipt", txnData)     // put data to session
	http.Redirect(w, r, "/receipt", http.StatusSeeOther) // redirect
}

// display receipt page
func (app *application) Receipt(w http.ResponseWriter, r *http.Request) {
	txn := app.Session.Get(r.Context(), "receipt").(TransactionData) // grap data from session
	data := make(map[string]any)
	data["txn"] = txn
	app.Session.Remove(r.Context(), "receipt")                   // remove data from session
	if err := app.renderTemplate(w, r, "receipt", &templateData{ // render receipt page
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

// display the page to subscribe bronze plan
func (app *application) BronzePlan(w http.ResponseWriter, r *http.Request) {
	widget, err := app.DB.GetWidget(2)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	data := make(map[string]any)
	data["widget"] = widget

	if err := app.renderTemplate(w, r, "bronze-plan", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

// display bronze plan(subscription) receipt page
func (app *application) BronzePlanReceipt(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "receipt-plan", nil); err != nil {
		app.errorLog.Println(err)
	}
}

type TransactionData struct {
	FirstName       string
	LastName        string
	Email           string
	PaymentIntentID string
	PaymentMethodID string
	PaymentAmount   int
	PaymentCurrency string
	LastFour        string
	ExpiryMonth     int
	ExpiryYear      int
	BankReturnCode  string
}

// get transaction data from posted form and stripe
func (app *application) GetTransactionData(r *http.Request) (TransactionData, error) {
	var txnData TransactionData
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return txnData, err
	}

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")

	amount, err := strconv.Atoi(paymentAmount)
	if err != nil {
		app.errorLog.Println(err)
		return txnData, err
	}

	card := cards.Card{
		Secret: app.config.stripe.secret,
		Key:    app.config.stripe.key,
	}

	pi, err := card.RetrievePaymentIntent(paymentIntent)
	if err != nil {
		app.errorLog.Println(err)
		return txnData, err
	}

	pm, err := card.GetPaymentMethod(paymentMethod)
	if err != nil {
		app.errorLog.Println(err)
		return txnData, err
	}

	lastFour := pm.Card.Last4
	expiryMonth := pm.Card.ExpMonth
	expiryYear := pm.Card.ExpYear

	txnData = TransactionData{
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		PaymentIntentID: paymentIntent,
		PaymentMethodID: paymentMethod,
		PaymentAmount:   amount,
		PaymentCurrency: paymentCurrency,
		LastFour:        lastFour,
		ExpiryMonth:     int(expiryMonth),
		ExpiryYear:      int(expiryYear),
		BankReturnCode:  pi.Charges.Data[0].ID,
	}

	return txnData, nil
}

// save customer information into database and return id
func (app *application) SaveCustomer(firstName, lastName, email string) (int, error) {
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	id, err := app.DB.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// save transaction information into database and return id
func (app *application) SaveTransaction(txn models.Transaction) (int, error) {
	id, err := app.DB.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// save order information into database and return id
func (app *application) SaveOrder(order models.Order) (int, error) {
	id, err := app.DB.InsertOrder(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// display login page
func (app *application) LoginPage(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "login", nil, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

// authenticate user and save userID to session
func (app *application) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	id, err := app.DB.Authenticate(email, password)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// clear the session
func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "forgot-password", nil); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ShowResetPassword(w http.ResponseWriter, r *http.Request) {
	reqURL := r.RequestURI
	testURL := fmt.Sprintf("%s%s", app.config.frontend, reqURL)

	signer := urlsigner.Signer{
		Secret: []byte(app.config.secretkey),
	}

	valid := signer.VerifyToken(testURL)

	// test
	if valid {
		w.Write([]byte("valid"))
	} else {
		w.Write([]byte("invalid"))
	}
}
