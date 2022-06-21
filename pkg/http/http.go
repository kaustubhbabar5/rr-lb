package http

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, body map[string]any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	//TODO handle error here
	_ = json.NewEncoder(w).Encode(body)

}
