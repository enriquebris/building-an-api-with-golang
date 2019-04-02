package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID      string `json:id"" valid:"required"`
	Name    string `json:"name" valid:"required"`
	Email   string `json:"email" valid:"email"`
	Address string `json:"address"`
}

type BasicResponse struct {
	Message string `json:"message"`
}

// initAPI starts the RESTful API endpoints
func initAPI(port string) {
	router := mux.NewRouter()

	//endpoints
	router.HandleFunc("/api/v1/ping", pingGET).Methods("GET")
	router.HandleFunc("/api/v1/status", statusGET).Methods("GET")

	router.HandleFunc(`/api/v1/user/{id:[a-zA-Z0-9\-]+}`, userGET).Methods("GET")
	router.HandleFunc(`/api/v1/user/`, userPOST).Methods("POST")

	http.ListenAndServe(port, router)
}

// outputJSON outputs json responses
func outputJSON(w http.ResponseWriter, httpCode int, jsonStruct interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	json.NewEncoder(w).Encode(jsonStruct)
}
