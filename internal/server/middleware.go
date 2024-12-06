package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type errorResponse struct {
	Err string `json:"error,omitempty"`
}

func middleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Enable CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// TODO: Make this configurable -- should we introduce gRPC metadata library & convert metadata from ctx into headers?
		// OK, I don't think this works with Vite dev hot reloading. Maybe the cache keeps refreshing?
		w.Header().Set("Cache-Control", "max-age=300, public")

		start := time.Now()
		handler(w, r)
		elapsed := time.Since(start).Microseconds()
		log.Printf("%s %s (%dμs)\n", r.Method, r.URL.Path, elapsed)
	}
}

func writeError(w http.ResponseWriter, err error) {
	outBytes, _ := json.Marshal(&errorResponse{Err: err.Error()})
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(outBytes)
}

func translateHandler(handler HandlerFunc) http.HandlerFunc {
	translated := func(w http.ResponseWriter, r *http.Request) {
		decode := func(req interface{}) error {
			err := json.NewDecoder(r.Body).Decode(req)
			if err == io.EOF {
				return nil
			}
			return err
		}

		out, err := handler(r.Context(), decode)
		if err != nil {
			writeError(w, err)
			return
		}

		outBytes, err := json.Marshal(out)
		if err != nil {
			writeError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(outBytes)
	}

	return middleware(translated)
}
