package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rivettax/unheic/unheicd/internal"
)

const (
	defaultReadTimeout  = 30
	defaultWriteTimeout = 30
	defaultIdleTimeout  = 60
	defaultPort         = 8080
)

func main() {
	port := getEnvAsInt("PORT", defaultPort)
	readTimeout := getEnvAsInt("READ_TIMEOUT", defaultReadTimeout)
	writeTimeout := getEnvAsInt("WRITE_TIMEOUT", defaultWriteTimeout)
	idleTimeout := getEnvAsInt("IDLE_TIMEOUT", defaultIdleTimeout)

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
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}

	log.Printf("Starting server on :%d", port)
	log.Fatal(server.ListenAndServe())
}

func handleHeifToJpeg(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "attachment; filename=converted.jpg")

	// Convert HEIF to JPEG
	err := internal.HeicToJPEG(r.Context(), w, r.Body)
	if err != nil {
		log.Printf("Error converting HEIF to JPEG: %v", err)
		http.Error(w, "Failed to convert image: "+err.Error(), http.StatusBadRequest)
		return
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid %s: %v", key, err)
	}

	return intValue
}
