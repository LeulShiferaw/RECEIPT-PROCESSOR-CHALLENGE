package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Constants for error message
var badRequest = "Please verify input."
var notFound = "No receipt found for that ID."

// Used for generating unique id for each request
// We can use UUID instead but I though this would be a simpler approach
var idCounter int

// Non persistent memory
var memory = make(map[string]Receipt)

// Item that is part of receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Receipt Datastructure for receiving json data
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

// Helper function for generating unique id
func generateID() int {
	idCounter++
	return idCounter
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", defaultHandler).Methods("GET")
	r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

	fmt.Println("Running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
