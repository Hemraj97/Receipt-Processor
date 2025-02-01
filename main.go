package main

import (
	"encoding/json"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ResponseID struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

var (
	receipts = make(map[string]Receipt)
	mutex    sync.Mutex
)

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	mutex.Lock()
	receipts[id] = receipt
	mutex.Unlock()

	json.NewEncoder(w).Encode(ResponseID{ID: id})
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/receipts/") : len(r.URL.Path)-len("/points")]
	mutex.Lock()
	receipt, exists := receipts[id]
	mutex.Unlock()

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points := calculatePoints(receipt)
	json.NewEncoder(w).Encode(PointsResponse{Points: points})
}

func calculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name
	retailerPoints := len(regexp.MustCompile(`[a-zA-Z0-9]`).FindAllString(receipt.Retailer, -1))
	points += retailerPoints

	// Rule 2: 50 points if total is a round dollar amount
	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == float64(int(total)) {
		points += 50
	}

	// Rule 3: 25 points if total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt
	pairPoints := (len(receipt.Items) / 2) * 5
	points += pairPoints

	// Rule 5: Trimmed item description length is a multiple of 3
	descriptionPoints := 0
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			itemPoints := int(math.Ceil(price * 0.2))
			descriptionPoints += itemPoints
		}
	}
	points += descriptionPoints

	// Rule 6: 6 points if the day of purchase is odd
	date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	oddDayPoints := 0
	if date.Day()%2 == 1 {
		oddDayPoints = 6
		points += oddDayPoints
	}

	// Rule 7: 10 points if purchase time is between 2:00 PM and 4:00 PM
	timeParts, _ := time.Parse("15:04", receipt.PurchaseTime)
	timeBonus := 0
	if timeParts.Hour() >= 14 && timeParts.Hour() < 16 {
		timeBonus = 10
		points += timeBonus
	}
	return points
}

func main() {
	http.HandleFunc("/receipts/process", processReceipt)
	http.HandleFunc("/receipts/", getPoints)
	http.ListenAndServe(":8080", nil)
}
