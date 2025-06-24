package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

func fetchTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request to Blockchair API
	resp, err := client.Get("https://api.blockchair.com/bitcoin/transactions")
	if err != nil {
		log.Printf("Error fetching from Blockchair API: %v", err)
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Parse Blockchair response
	var blockchairResp BlockchairResponse
	if err := json.Unmarshal(body, &blockchairResp); err != nil {
		log.Printf("Error parsing Blockchair response: %v", err)
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	if len(blockchairResp.Data) == 0 {
		log.Printf("No transactions found in response")
		http.Error(w, "No transactions found", http.StatusNotFound)
		return
	}

	// Create response
	response := TransactionsResponse{
		Success: true,
		Data:    blockchairResp.Data,
	}

	// Send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
