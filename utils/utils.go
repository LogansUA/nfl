package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
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

func RandToken(len int) string {
	b := make([]byte, len)

	rand.Read(b)

	return fmt.Sprintf("%x", b)
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
