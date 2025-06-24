package main

// BlockchairResponse represents the response structure from Blockchair API
// Update: data is a slice of transactions, not a map
// context is a map
type BlockchairResponse struct {
	Data    []Transaction          `json:"data"`
	Context map[string]interface{} `json:"context"`
}

// Transaction represents a Bitcoin transaction
type Transaction struct {
	BlockID        int64   `json:"block_id"`
	ID             int64   `json:"id"`
	Hash           string  `json:"hash"`
	Date           string  `json:"date"`
	Time           string  `json:"time"`
	Size           int     `json:"size"`
	Weight         int     `json:"weight"`
	Version        int     `json:"version"`
	LockTime       int64   `json:"lock_time"`
	IsCoinbase     bool    `json:"is_coinbase"`
	HasWitness     bool    `json:"has_witness"`
	InputCount     int     `json:"input_count"`
	OutputCount    int     `json:"output_count"`
	InputTotal     int64   `json:"input_total"`
	InputTotalUsd  float64 `json:"input_total_usd"`
	OutputTotal    int64   `json:"output_total"`
	OutputTotalUsd float64 `json:"output_total_usd"`
	Fee            int64   `json:"fee"`
	FeeUsd         float64 `json:"fee_usd"`
	FeePerKb       float64 `json:"fee_per_kb"`
	FeePerKbUsd    float64 `json:"fee_per_kb_usd"`
	FeePerKwu      float64 `json:"fee_per_kwu"`
	FeePerKwuUsd   float64 `json:"fee_per_kwu_usd"`
	CddTotal       float64 `json:"cdd_total"`
}

// TransactionsResponse represents the response for our transactions endpoint
type TransactionsResponse struct {
	Success bool          `json:"success"`
	Data    []Transaction `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}
