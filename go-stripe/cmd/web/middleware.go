package main

import "net/http"

// set up middleware that loads and saves session automatically
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
