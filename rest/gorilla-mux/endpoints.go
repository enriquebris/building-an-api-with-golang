package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

// pingGet handles [GET] /api/v1/ping
func pingGET(w http.ResponseWriter, req *http.Request) {
	outputJSON(w, http.StatusOK, BasicResponse{"pong"})
}

// statusGET handles [GET] /api/v1/status
// It checks the context and acts in a different way if it was canceled during the execution.
func statusGET(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	totalStatuses := 5
	statusChannel := make(chan string, totalStatuses)
	for i := 0; i < totalStatuses; i++ {
		go getStatusHelper(ctx, fmt.Sprintf("job-%v", i), statusChannel)
	}

	total := 0
	for total < totalStatuses {
		select {
		case status := <-statusChannel:
			fmt.Println(status)
			total++
		}
	}

	fmt.Println("done")
}

// getStatusHelper returns a dummy status, but checks the context
func getStatusHelper(ctx context.Context, name string, ch chan<- string) {
	// do the job
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	randomNumber := r.Intn(10)
	sleeptime := randomNumber + 10
	fmt.Println(name, "Starting sleep for", sleeptime, "s")
	time.Sleep(time.Duration(sleeptime) * time.Second)

	select {
	case <-ctx.Done():
		// return a different message if the context was canceled
		ch <- fmt.Sprintf("Status canceled for '%v'", name)
		return
	default:
		ch <- fmt.Sprintf("Status for '%v'", name)
		return
	}
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
