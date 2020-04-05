package handler

import (
	"encoding/json"
	"log"
)

// errorResponse represents http response when error occurred
type errorResponse struct {
	Message string `json:"message"`
}

func logWithJSON(prefix string, data interface{}) {
	contents, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal data for log: %v", err)
		return
	}
	log.Println(prefix + ": " + string(contents))
}
