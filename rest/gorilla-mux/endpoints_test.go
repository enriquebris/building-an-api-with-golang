package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestPingGET(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/api/v1/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(pingGET)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expectedBody, _ := json.Marshal(BasicResponse{"pong"})
	currentBody := strings.TrimSpace(rr.Body.String())
	if currentBody != string(expectedBody) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expectedBody))
	}
}

func TestUserGET(t *testing.T) {
	name := "testingID"

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/v1/user/%v", name), nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Use a router from mux in order to add the params/vars to the context
	router := mux.NewRouter()
	router.HandleFunc(`/api/v1/user/{id:[a-zA-Z0-9\-]+}`, userGET).Methods("GET")
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expectedBody, _ := json.Marshal(BasicResponse{fmt.Sprintf("Hello %v!", name)})
	currentBody := strings.TrimSpace(rr.Body.String())
	if currentBody != string(expectedBody) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expectedBody))
	}
}

func TestUserPOSTUsingValidData(t *testing.T) {
	// build the test payload
	user := User{
		ID:      "newID",
		Name:    "John Doe",
		Address: "P.O. BOX 12345",
		Email:   "john@doe.xyz",
	}
	payloadJSON, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/user/", bytes.NewReader(payloadJSON))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userPOST)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expectedBody, _ := json.Marshal(BasicResponse{fmt.Sprintf("User '%v' (id: %v) was successfully added.", user.Name, user.ID)})
	currentBody := strings.TrimSpace(rr.Body.String())
	if currentBody != string(expectedBody) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), string(expectedBody))
	}
}
