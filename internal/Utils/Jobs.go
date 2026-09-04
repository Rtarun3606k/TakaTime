package utils

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	gogist "github.com/Rtarun3606k/TakaTime/internal/GoGist"
)

func HandleImageJob(name, path, token, repo string, generator func() (image.Image, error)) {
	fmt.Printf("Processing %s...\n", name)

	// 1. Generate Image
	img, err := generator()
	if err != nil {
		log.Printf("Gen Error (%s): %v\n", name, err)
		return
	}

	//debugging image generation
	// SaveImage(name+".png", img)

	// 2. Format Config (Using your utils package)
	cfg, err := FormmatUpload(token, repo, path, "main", "Update "+name)
	if err != nil {
		log.Printf("Config Error (%s): %v\n", name, err)
		return
	}

	// 3. Upload with FRESH Timeout (Critical for loops!)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // Cancels this specific context when function exits

	if err := gogist.UploadImageToGitHub(ctx, img, cfg); err != nil {
		log.Printf("Upload Error (%s): %v\n", name, err)
	} else {
		fmt.Printf("Uploaded: %s\n", path)
	}
}

func SaveImage(filename string, img image.Image) error {
	// 1. Create the file
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// 2. Encode the image as PNG
	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	fmt.Printf("Saved debug image: %s\n", filename)
	return nil
}
