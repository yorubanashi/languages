package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// writeJSON will attempt to convert 'data' into JSON and write it to http.ResponseWriter.
//
// If unsuccessful, it will instead write a 500 Internal Server Error.
func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	out, err := json.Marshal(data)
	if err != nil {
		errorStr := fmt.Sprintf("error marshalling JSON: %s", err.Error())
		w.WriteHeader(500)
		w.Write([]byte(errorStr))
	} else {
		w.WriteHeader(statusCode)
		w.Write(out)
	}
}
