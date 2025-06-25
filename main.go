package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := "8080"
	router := mux.NewRouter()

	// Apply metrics middleware to all routes
	router.Use(metricsMiddleware)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Register endpoints
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/transactions", fetchTransactionsHandler).Methods("GET")
	router.Handle("/metrics", promhttp.Handler())

	// Start collecting runtime metrics
	go collectRuntimeMetrics()

	// create stop signal listener: SIGINT is sent from CTRL+C, SIGTERM is sent from K8s
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("Error starting server: %s", err)
		}
	}()
	log.Printf("Server listening on port %s\n", port)
	log.Printf("Metrics available at http://localhost:%s/metrics\n", port)
	log.Printf("Health check available at http://localhost:%s/health\n", port)

	<-stop

	log.Printf("Server beginning graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)

	// release connections
	defer func() {
		// _ = log.Sync()
		fmt.Println("Released connections")
		cancel()
	}()

	// shutdown server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %s", err)

	}
	fmt.Println("Server shutdown complete")
}
