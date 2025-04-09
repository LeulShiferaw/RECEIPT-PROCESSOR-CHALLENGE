package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var badRequest = "Please verify input."
var notFound = "No receipt found for that ID."

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", defaultHandler).Methods("GET")
	r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

	fmt.Println("Running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Fetch rewards rocks!"))
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	fmt.Printf("Received receipt: %+v\n", receipt) //For Debugging purposes
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}
