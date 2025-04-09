package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gorilla/mux"
)

var badRequest = "Please verify input."
var notFound = "No receipt found for that ID."

var idCounter int

var memory = make(map[string]Receipt)

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
	id := strconv.Itoa(generateID())               //Create new Id
	memory[id] = receipt                           //Insert to memory

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("Calculating for ID: %s\n", id)

	receipt, found := memory[id]
	if !found {
		http.Error(w, notFound, http.StatusNotFound)
		return
	}

	fmt.Printf("Calculating for receipt: %+v\n", receipt)

	points := calcPoints(receipt)
	fmt.Printf("Points calculated: %f\n", points)

	if points < 0 {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]int{"points": int(points)})
}

func isInteger(f float64) bool {
	return float64(int64(f)) == f
}

func calcPoints(receipt Receipt) float64 {
	//Points for every alphanumeric character in receipt.Retailer
	var points float64
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points++
		}
	}

	fmt.Printf("Added %f at the start\n", points) //Debugging help

	//Get total from string
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		fmt.Println("Error with parsing total!")
		return -1.0
	}

	//Check if total is a round dollar amount
	if isInteger(total) {
		fmt.Println("Added 50 for round dollar total")
		points += 50
	}

	//Check if total is multiple of 0.25
	if isInteger(total * 4) {
		fmt.Println("Added 25 for multiple of 4 integer")
		points += 25
	}

	fmt.Printf("Adding %f points based on length of items\n", 5*float64(int64(len(receipt.Items)/2)))
	points += 5 * float64(int64(len(receipt.Items)/2)) //5 points for every two items

	//Add ceil(price*0.2) for descriptions that have trimmed length of multiple 3
	for _, item := range receipt.Items {
		trimmedItem := strings.TrimSpace(item.ShortDescription)
		if len(trimmedItem)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				fmt.Println("Error with parsing price!")
				return -6.0
			}
			fmt.Printf("Added %f for item: %s\n", math.Ceil(price*0.2), trimmedItem)
			points += math.Ceil(price * 0.2)
		}
	}

	//AI help
	if total > 10 {
		fmt.Println("Added 5 for AI help")
		points += 5
	}

	//6 for odd day
	date := strings.Split(receipt.PurchaseDate, "-")
	if len(date) >= 3 {
		day, err := strconv.Atoi(date[2])
		if err != nil {
			fmt.Println("Error with purchaseDate!")
			return -2.0
		}

		if (day % 2) != 0 {
			fmt.Println("Added 6 for odd day")
			points += 6
		}
	}

	//Get time information
	time := strings.Split(receipt.PurchaseTime, ":")
	if len(time) != 2 {
		fmt.Println("Error with len(time)")
		return -3.0
	}

	hour, err := strconv.Atoi(time[0])
	if err != nil {
		fmt.Println("Error with parsing hour from purchase time!")
		return -4.0
	}

	min, err := strconv.Atoi(time[1])
	if err != nil {
		fmt.Println("Error with parsing minute from purchase time!")
		return -5.0
	}

	//Check if 24 hour format is between 2:00pm and 4:00pm
	if ((hour > 14) || (hour == 14 && min > 0)) && (hour < 16) {
		fmt.Println("Added 10 for hour between 2pm and 4pm")
		points += 10
	}

	return points
}
