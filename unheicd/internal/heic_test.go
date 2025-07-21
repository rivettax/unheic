package internal

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"testing"
)

// Test for HEIC to JPEG conversion
func TestHeicToJPEG(t *testing.T) {
	// Open the sample HEIC file for testing
	heicFile, err := os.Open("../testdata/test.heic")
	if err != nil {
		t.Fatalf("opening test file: %v", err)
	}
	defer heicFile.Close()

	// Create a temporary file to store the converted JPEG output
	tempJpegFile, err := os.CreateTemp("", "*.jpg")
	if err != nil {
		t.Fatalf("creating tempJpegFile: %v", err)
	}
	defer tempJpegFile.Close()

	// Call the function to convert HEIC to JPEG
	err = HeicToJPEG(context.Background(), tempJpegFile, heicFile)
	if err != nil {
		t.Fatalf("converting HEIC to JPEG: %v", err)
	}

	// Reset file pointer to the beginning of the temp JPEG file for hashing
	_, err = tempJpegFile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("seeking to start of tempJpegFile: %v", err)
	}
	// Compute SHA256 hash of the generated JPEG file
	hasher := sha256.New()
	if _, err := io.Copy(hasher, tempJpegFile); err != nil {
		t.Fatalf("hashing tempJpegFile: %v", err)
	}
	tempJpegFileHash := hasher.Sum(nil)

	// Open the reference JPEG file for comparison
	jpegFile, err := os.Open("../testdata/test.jpg")
	if err != nil {
		t.Fatalf("opening reference file: %v", err)
	}
	defer jpegFile.Close()

	// Compute SHA256 hash of the reference JPEG file
	hasher2 := sha256.New()
	if _, err := io.Copy(hasher2, jpegFile); err != nil {
		t.Fatalf("hashing reference file: %v", err)
	}
	jpegFileHash := hasher2.Sum(nil)

	// Print out the hashes for debugging purposes
	fmt.Printf("tempJpegFileHash: %x\n", tempJpegFileHash)
	fmt.Printf("jpegFileHash: %x\n", jpegFileHash)

	// Compare the two hashes to ensure the conversion output matches the reference
	if !bytes.Equal(tempJpegFileHash, jpegFileHash) {
		t.Fatalf("output files do not match")
	}

}
