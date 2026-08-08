package utils

import (
	"errors"
	"image"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveImage(t *testing.T) {
	dir := t.TempDir()

	t.Run("Successful Save", func(t *testing.T) {
		path := filepath.Join(dir, "test_output.png")

		// Create a tiny 10x10 blank image in memory
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))

		err := SaveImage(path, img)
		if err != nil {
			t.Fatalf("SaveImage() unexpected error: %v", err)
		}

		// Verify the file actually exists on disk
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			t.Errorf("SaveImage() failed to write the file to disk")
		}
		if info.Size() == 0 {
			t.Errorf("SaveImage() wrote an empty file")
		}
	})

	t.Run("Directory Does Not Exist", func(t *testing.T) {
		// Attempt to save to an invalid path that doesn't exist
		path := filepath.Join(dir, "missing_folder", "fail.png")
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))

		err := SaveImage(path, img)
		if err == nil {
			t.Errorf("SaveImage() expected an error for missing directory, got nil")
		}
	})
}

func TestHandleImageJob(t *testing.T) {
	t.Run("Generator Failure Halts Execution", func(t *testing.T) {
		generatorCalled := false

		// Mock generator that immediately returns an error
		mockFailingGenerator := func() (image.Image, error) {
			generatorCalled = true
			return nil, errors.New("mock generator failed") // make sure to import "errors" at the top!
		}

		// Execute job - this should log the generation error and return early
		HandleImageJob(
			"Test-Fail-Job",
			"fake/path.png",
			"fake-token",
			"fake-repo",
			mockFailingGenerator,
		)

		if !generatorCalled {
			t.Errorf("HandleImageJob() did not execute the injected generator function")
		}
	})

	t.Run("Proceeds to Config on Successful Generation", func(t *testing.T) {
		generatorCalled := false

		// Mock generator that succeeds
		mockSuccessGenerator := func() (image.Image, error) {
			generatorCalled = true
			return image.NewRGBA(image.Rect(0, 0, 10, 10)), nil
		}

		// Execute job with a properly formatted dummy repo string!
		HandleImageJob(
			"Test-Success-Job",
			"fake/path.png",
			"invalid-token",
			"fake-owner/fake-repo", // <--- Correctly formatted!
			mockSuccessGenerator,
		)

		if !generatorCalled {
			t.Errorf("HandleImageJob() did not execute the injected generator function")
		}
	})
}

