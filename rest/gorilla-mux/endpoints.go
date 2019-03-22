package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// pingGet handles [GET] /api/v1/ping
func pingGET(w http.ResponseWriter, req *http.Request) {
	outputJSON(w, http.StatusOK, BasicResponse{"pong"})
}

// userGET handles [GET] /api/v1/user/{id:[a-zA-Z0-9\-]+}
func userGET(w http.ResponseWriter, req *http.Request) {
	// read url arguments
	vars := mux.Vars(req)
	name := vars["id"]

	outputJSON(w, http.StatusOK, BasicResponse{fmt.Sprintf("Hello %v!", name)})
}

// userPOST handles [POST] /api/v1/user/
func userPOST(w http.ResponseWriter, req *http.Request) {
	var user User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		outputJSON(w, http.StatusBadRequest, BasicResponse{err.Error()})
		return
	}

	// validate the payload
	if okValidation, err := govalidator.ValidateStruct(user); err != nil {
		if okValidation {
			outputJSON(w, http.StatusInternalServerError, BasicResponse{err.Error()})
		}

		outputJSON(w, http.StatusBadRequest, BasicResponse{err.Error()})
		return
	}

	// do the magic here ...

	user.ID = "newID"
	outputJSON(w, http.StatusOK, BasicResponse{fmt.Sprintf("User '%v' (id: %v) was successfully added.", user.Name, user.ID)})
}
