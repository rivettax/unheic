package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rivettax/unheic/internal"
)

func main() {
	// Create a new HTTP server
	mux := http.NewServeMux()

	// Register the HEIF to JPEG conversion endpoint
	mux.HandleFunc("POST /convert", handleHeifToJpeg)

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Configure server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on :8080")
	log.Fatal(server.ListenAndServe())
}

func handleHeifToJpeg(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "attachment; filename=converted.jpg")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Convert HEIF to JPEG
	err := internal.HeicToJPEG(ctx, w, r.Body)
	if err != nil {
		log.Printf("Error converting HEIF to JPEG: %v", err)
		http.Error(w, "Failed to convert image: "+err.Error(), http.StatusBadRequest)
		return
	}
}
