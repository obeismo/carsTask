package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func newResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}
