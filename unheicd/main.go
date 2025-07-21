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
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// configuration options are provided as environment variables with reasonable defaults
	port := getEnvAsInt("PORT", defaultPort)
	readTimeout := getEnvAsInt("READ_TIMEOUT", defaultReadTimeout)
	writeTimeout := getEnvAsInt("WRITE_TIMEOUT", defaultWriteTimeout)
	idleTimeout := getEnvAsInt("IDLE_TIMEOUT", defaultIdleTimeout)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /convert", handleHeifToJpeg)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}

	log.Printf("Starting server on :%d", port)

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func handleHeifToJpeg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "attachment; filename=converted.jpg")

	err := internal.HeicToJPEG(r.Context(), w, r.Body)
	if err != nil {
		switch err.(type) {
		case internal.DecodeError:
			http.Error(w, "decoding HEIF image: "+err.Error(), http.StatusBadRequest)
		case internal.EncodeError:
			http.Error(w, "encoding JPEG image: "+err.Error(), http.StatusInternalServerError)
		default:
			http.Error(w, "converting HEIF to JPEG: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}

	return intValue
}
