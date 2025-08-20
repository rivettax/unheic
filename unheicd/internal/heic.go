package internal

import (
	"context"
	"fmt"
	"image/jpeg"
	"io"

	"github.com/jdeng/goheif"
)

type DecodeError struct{ Err error }

func (e DecodeError) Error() string { return e.Err.Error() }

type EncodeError struct{ Err error }

func (e EncodeError) Error() string { return e.Err.Error() }

func HeicToJPEG(ctx context.Context, w io.Writer, r io.Reader) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = DecodeError{fmt.Errorf("panic in HEIC decoder: %v", r)}
		}
	}()

	img, err := goheif.Decode(r)
	if err != nil {
		return DecodeError{err}
	}

	err = jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
	if err != nil {
		return EncodeError{err}
	}

	return nil
}
