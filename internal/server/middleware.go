package server

import (
	"log"
	"net/http"
)

func middleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		handler(w, r)
		log.Printf("%s %s\n", r.Method, r.URL.Path)
	}
}
