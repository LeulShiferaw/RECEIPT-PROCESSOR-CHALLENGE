# RECEIPT-PROCESSOR-CHALLENGE
My Solution to receipt processor challenge for fetch rewards

## Steps to run  

# Installing dependencies
1. go install github.com/gorilla/mux  
2. go get github.com/gorilla/mux

# Exectuing program
1. go run main.go handlers.go points.go 

# Steps to run tests
1. go test -v

## Additional Information
- This service will run on http://localhost:8080
- It can be tested using curl commands
- Example for /receipts/process:
    curl -X POST http://localhost:8080/receipts/process   -H "Content-Type: application/json"   -d '{
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
        {
        "shortDescription": "Mountain Dew 12PK",
        "price": "6.49"
        },{
        "shortDescription": "Emils Cheese Pizza",
        "price": "12.25"
        },{
        "shortDescription": "Knorr Creamy Chicken",
        "price": "1.26"
        },{
        "shortDescription": "Doritos Nacho Cheese",
        "price": "3.35"
        },{
        "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
        "price": "12.00"
        }
    ],
    "total": "35.35"
    }'
- Example for /receipts/{id}/points:
    curl http://localhost:8080/receipts/1/points
- I decided to use a simple id just "1", "2", "3" ...
  I did this to make the project simpler and include as little dependencies as possible
- The mux dependency is used for ensuring that the API endpoints specifically use POST or GET.