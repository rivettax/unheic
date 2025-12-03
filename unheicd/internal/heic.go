package internal

import (
	"context"
	"fmt"
	"image/jpeg"
	"io"

	"github.com/strukturag/libheif/go/heif"
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

	data, err := io.ReadAll(r)
	if err != nil {
		return DecodeError{err}
	}

	context, err := heif.NewContext()
	if err != nil {
		return DecodeError{err}
	}

	err = context.ReadFromMemory(data)
	if err != nil {
		return DecodeError{err}
	}

	handle, err := context.GetPrimaryImageHandle()
	if err != nil {
		return DecodeError{err}
	}

	img, err := handle.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
	if err != nil {
		return DecodeError{err}
	}

	goImg, err := img.GetImage()
	if err != nil {
		return DecodeError{err}
	}

	err = jpeg.Encode(w, goImg, &jpeg.Options{Quality: 90})
	if err != nil {
		return EncodeError{err}
	}

	return nil
}
