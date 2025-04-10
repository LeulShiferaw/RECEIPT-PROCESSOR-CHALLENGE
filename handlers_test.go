package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")
	return r
}

func TestProcessReceiptHandler(t *testing.T) {
	router := setupRouter()

	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "15:01",
		Total:        "20.00",
		Items: []Item{
			{ShortDescription: "Pepsi", Price: "2.00"},
			{ShortDescription: "Chips", Price: "3.00"},
		},
	}

	payload, _ := json.Marshal(receipt)

	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", resp.Code)
	}

	var result map[string]string
	json.Unmarshal(resp.Body.Bytes(), &result)

	if _, exists := result["id"]; !exists {
		t.Errorf("Expected an ID in response")
	}
}

func TestGetPointsHandler(t *testing.T) {
	router := setupRouter()

	// First, create a receipt
	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "15:01",
		Total:        "20.00",
		Items: []Item{
			{ShortDescription: "Pepsi", Price: "2.00"},
			{ShortDescription: "Chips", Price: "3.00"},
		},
	}
	payload, _ := json.Marshal(receipt)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var result map[string]string
	json.Unmarshal(resp.Body.Bytes(), &result)
	id := result["id"]

	// Now fetch the points
	req2, _ := http.NewRequest("GET", "/receipts/"+id+"/points", nil)
	resp2 := httptest.NewRecorder()
	router.ServeHTTP(resp2, req2)

	if resp2.Code != http.StatusOK {
		t.Errorf("Expected 200 OK for points endpoint, got %d", resp2.Code)
	}

	var points map[string]int
	json.Unmarshal(resp2.Body.Bytes(), &points)
	if _, exists := points["points"]; !exists {
		t.Errorf("Expected points in response")
	}
}
