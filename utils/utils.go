package utils

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
)

func RenderError(w http.ResponseWriter, message string, statusCode int) {
	code := http.StatusBadRequest

	if statusCode != 0 {
		code = statusCode
	}

	w.WriteHeader(code)

	_, err := w.Write([]byte(message))

	log.Fatal(err)
}

func RandToken() string {
	return uuid.Must(uuid.NewV4()).String()
}

func JsonResponse(w http.ResponseWriter, model interface{}) {
	js, err := json.Marshal(model)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
