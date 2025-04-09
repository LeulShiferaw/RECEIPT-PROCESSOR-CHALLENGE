package main

import "testing"

func TestCalcPoints_RetailerName(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Target123",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "10:00",
		Total:        "1.1",
		Items:        []Item{},
	}

	points := calcPoints(receipt)

	if points != 9 { //5 for AI help
		t.Errorf("Expected 9 points, got %f", points)
	}
}

func TestCalcPoints_RoundDollar(t *testing.T) {
	//With round dollar = 50
	receipt := Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "10:00",
		Total:        "1.0",
		Items:        []Item{},
	}
	points := calcPoints(receipt)

	//Without round dollar = 0
	receipt = Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "10:00",
		Total:        "1.3",
		Items:        []Item{},
	}
	points2 := calcPoints(receipt)

	if points != 75 { //If it is a round dollar it is also a multiple of 0.25 (+25)
		t.Errorf("Expected 75 points, got %f", points)
	}

	if points2 != 0 {
		t.Errorf("Expected 0 points, got %f", points2)
	}
}

func TestCalcPoints_MultipleQuarter(t *testing.T) {
	receipt := Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "10:00",
		Total:        "1.1",
		Items: []Item{
			{ShortDescription: " THE", Price: "3.0"}, //Space at the beginning and three characters; should be incremented by ceil(3*0.2)=1
		},
	}
	points := calcPoints(receipt)

	if points != 1 {
		t.Errorf("Expected 1, got %f", points)
	}

	receipt = Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "10:00",
		Total:        "1.1",
		Items: []Item{
			{ShortDescription: " THEA", Price: "3.0"}, //Four characters so no increments
		},
	}
	points = calcPoints(receipt)

	if points != 0 {
		t.Errorf("Expected 0, got %f", points)
	}
}

func TestCalcPoints_AIHelp(t *testing.T) {
	receipt := Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "10:00",
		Total:        "11.1", //Must be above 10; should be incremented by 5
		Items:        []Item{},
	}
	points := calcPoints(receipt)

	if points != 5 {
		t.Errorf("Expected 5, got %f", points)
	}
}

func TestCalcPoints_OddDay(t *testing.T) {
	receipt := Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "10:00",
		Total:        "1.1",
		Items:        []Item{},
	}
	points := calcPoints(receipt)

	if points != 6 {
		t.Errorf("Expected 6, got %f", points)
	}
}

func TestCalcPoints_TimeBetween(t *testing.T) {
	receipt := Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:00", //Should be before 14:00 so no increment
		Total:        "1.1",
		Items:        []Item{},
	}
	points := calcPoints(receipt)

	if points != 0 {
		t.Errorf("Expected 0, got %f", points)
	}

	receipt = Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "14:00", //Should be at 14 still no increment
		Total:        "1.1",
		Items:        []Item{},
	}
	points = calcPoints(receipt)

	if points != 0 {
		t.Errorf("Expected 0, got %f", points)
	}

	receipt = Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "14:30", //Should be between, so increment by 10
		Total:        "1.1",
		Items:        []Item{},
	}
	points = calcPoints(receipt)

	if points != 10 {
		t.Errorf("Expected 10, got %f", points)
	}

	receipt = Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "15:30", //Should be between, so increment by 10
		Total:        "1.1",
		Items:        []Item{},
	}
	points = calcPoints(receipt)

	if points != 10 {
		t.Errorf("Expected 10, got %f", points)
	}

	receipt = Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "16:30", //Should be at 16, so no increment by 10
		Total:        "1.1",
		Items:        []Item{},
	}
	points = calcPoints(receipt)

	if points != 0 {
		t.Errorf("Expected 0, got %f", points)
	}

	receipt = Receipt{
		Retailer:     "",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "16:30", //Should after 16, so no increment by 10
		Total:        "1.1",
		Items:        []Item{},
	}
	points = calcPoints(receipt)

	if points != 0 {
		t.Errorf("Expected 0, got %f", points)
	}
}

// Test the examples given in the problem
func TestCalcPoints_WholeTests(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-01",
		PurchaseTime: "13:01",
		Total:        "35.35",
		Items: []Item{
			{ShortDescription: "Mountain Dew 12PK", Price: "6.49"},
			{ShortDescription: "Emils Cheese Pizza", Price: "12.25"},
			{ShortDescription: "Knorr Creamy Chicken", Price: "1.26"},
			{ShortDescription: "Doritos Nacho Cheese", Price: "3.35"},
			{ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ", Price: "12.00"},
		},
	}
	points := calcPoints(receipt)

	if points-5 != 28 { //Subract 5 for AI help
		t.Errorf("Expected 28, got %f", points)
	}

	receipt = Receipt{
		Retailer:     "M&M Corner Market",
		PurchaseDate: "2022-03-20",
		PurchaseTime: "14:33",
		Total:        "9.00",
		Items: []Item{
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
			{ShortDescription: "Gatorade", Price: "2.25"},
		},
	}
	points = calcPoints(receipt)

	if points != 109 { //No AI help because total is less than 10
		t.Errorf("Expected 28, got %f", points)
	}
}
