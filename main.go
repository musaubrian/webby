package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func eventHandler(w http.ResponseWriter, r *http.Request) {
	var event Event
	var count int64

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)

	fData := &FileSave{
		Id:    count,
		Event: event,
	}
	res, err := json.Marshal(fData)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	if err := os.WriteFile("responses.json", res, 0644); err != nil {
		http.Error(w, "Failed to save to file", http.StatusInternalServerError)
		return
	}
	log.Println("Saved to file")
	count += 1
}

// Show the last webhook you received
func prevWebhooks(w http.ResponseWriter, r *http.Request) {
	var all FileSave
	data, err := os.ReadFile("responses.json")
	if err != nil {
		http.Error(w, "Failed to fetch all webhook responses", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &all)
	if err != nil {
		http.Error(w, "Failed to load all webhooks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(all)
}

func main() {
	http.HandleFunc("/", eventHandler)
	http.HandleFunc("/prev", prevWebhooks)
	log.Println("Server is running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
