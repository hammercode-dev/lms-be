package utils

import (
	"encoding/json"
	"net/http"

	"github.com/hammer-code/lms-be/domain"
)

func Response(data domain.HttpResponse, w http.ResponseWriter) {
	// Marshal the data to JSON
	response, err := json.Marshal(data)
	if err != nil {
		// Handle error in case of failed JSON marshaling
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	// Set Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Set the status code as provided
	w.WriteHeader(data.Code)
	// Write the JSON data
	w.Write(response)
}
