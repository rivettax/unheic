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

func TestHeicToJPEG(t *testing.T) {
	tests := []struct {
		name     string
		heicFile string
		jpegFile string
	}{
		{name: "test.heic", heicFile: "../../testdata/test.heic", jpegFile: "../../testdata/test.jpg"},
		{name: "camel.heic", heicFile: "../../testdata/camel.heic", jpegFile: "../../testdata/camel.jpg"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			heicFile, err := os.Open(test.heicFile)
			if err != nil {
				t.Fatalf("opening test file: %v", err)
			}
			defer heicFile.Close()

			tempJpegFile, err := os.CreateTemp("", "*.jpg")
			if err != nil {
				t.Fatalf("creating tempJpegFile: %v", err)
			}
			defer tempJpegFile.Close()

			err = HeicToJPEG(context.Background(), tempJpegFile, heicFile)
			if err != nil {
				t.Fatalf("converting HEIC to JPEG: %v", err)
			}

			_, err = tempJpegFile.Seek(0, io.SeekStart)
			if err != nil {
				t.Fatalf("seeking to start of tempJpegFile: %v", err)
			}
			hasher := sha256.New()
			if _, err := io.Copy(hasher, tempJpegFile); err != nil {
				t.Fatalf("hashing tempJpegFile: %v", err)
			}
			tempJpegFileHash := hasher.Sum(nil)

			jpegFile, err := os.Open(test.jpegFile)
			if err != nil {
				t.Fatalf("opening reference file: %v", err)
			}
			defer jpegFile.Close()

			hasher2 := sha256.New()
			if _, err := io.Copy(hasher2, jpegFile); err != nil {
				t.Fatalf("hashing reference file: %v", err)
			}
			jpegFileHash := hasher2.Sum(nil)

			fmt.Printf("tempJpegFileHash: %x\n", tempJpegFileHash)
			fmt.Printf("jpegFileHash: %x\n", jpegFileHash)

			if !bytes.Equal(tempJpegFileHash, jpegFileHash) {
				t.Fatalf("output files do not match")
			}
		})
	}
}
