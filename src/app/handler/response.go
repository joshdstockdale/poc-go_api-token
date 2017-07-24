package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleJsonResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	log.Println("Response: ", data)
	js, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SetContentTypeJson(w)
	w.Write(js)
}
