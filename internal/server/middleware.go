package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

		handler(w, r)
		log.Printf("%s %s\n", r.Method, r.URL.Path)
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
		return
	}

	return middleware(translated)
}
