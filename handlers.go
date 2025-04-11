package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Initial page that is not required
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Fetch rewards rocks!\nUsage:\n1. /receipts/process as POST with json data to get id\n2. /receipts/{id}/points to get the points associated with id"))
}

// The POST api endpoint for initial processing
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	//fmt.Printf("Received receipt: %+v\n", receipt) //For Debugging purposes
	id := strconv.Itoa(generateID()) //Create new Id
	memory[id] = receipt             //Insert to memory

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	//Get id
	vars := mux.Vars(r)
	id := vars["id"]

	//fmt.Printf("Calculating for ID: %s\n", id) //Debugging purposes

	//Check if the id has been assigned before and get the receipt that corresponds to it
	receipt, found := memory[id]
	if !found {
		http.Error(w, notFound, http.StatusNotFound)
		return
	}

	//fmt.Printf("Calculating for receipt: %+v\n", receipt)

	points := calcPoints(receipt)
	//fmt.Printf("Points calculated: %f\n", points)

	//calcPoints returns < 0 if there is an error
	if points < 0 {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]int{"points": int(points)})
}
