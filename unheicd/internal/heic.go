package internal

import (
	"context"
	"image/jpeg"
	"io"

	"github.com/jdeng/goheif"
)

type DecodeError struct{ Err error }

func (e DecodeError) Error() string { return e.Err.Error() }

type EncodeError struct{ Err error }

func (e EncodeError) Error() string { return e.Err.Error() }

func HeicToJPEG(ctx context.Context, w io.Writer, r io.Reader) error {
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
