package internal

import (
	"context"
	"fmt"
	"image/jpeg"
	"io"

	"github.com/jdeng/goheif"
)

func HeicToJPEG(ctx context.Context, w io.Writer, r io.Reader) error {
	// Decode the HEIC image
	img, err := goheif.Decode(r)
	if err != nil {
		return fmt.Errorf("decoding HEIC image: %w", err)
	}

	// Set JPEG encoding options
	options := jpeg.Options{
		Quality: 90,
	}

	// Encode as JPEG
	err = jpeg.Encode(w, img, &options)
	if err != nil {
		return fmt.Errorf("encoding JPEG image: %w", err)
	}

	return nil
}
