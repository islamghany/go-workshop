package main

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, data interface{}) {
	jsonData, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func SayHello(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, r.URL.Path+" hello")
}
func SayGoodbye(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, r.URL.Path+" goodbye")
}
