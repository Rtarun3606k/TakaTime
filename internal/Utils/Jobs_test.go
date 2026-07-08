package utils

import (
	"errors"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveImage(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})

	file := filepath.Join(t.TempDir(), "test.png")

	err := SaveImage(file, img)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(file)
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}

	if info.Size() == 0 {
		t.Fatal("image file is empty")
	}
}

func TestSaveImageInvalidPath(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))

	err := SaveImage("/invalid/path/test.png", img)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHandleImageJobGeneratorError(t *testing.T) {
	called := false

	HandleImageJob(
		"test-image",
		"test.png",
		"token",
		"repo",
		func() (image.Image, error) {
			called = true
			return nil, errors.New("generation failed")
		},
	)

	if !called {
		t.Fatal("generator function was not called")
	}
}
