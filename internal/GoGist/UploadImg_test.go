package gogist

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"golang.org/x/oauth2"
)

// mockTransport implements http.RoundTripper to intercept GitHub API calls natively
type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestUploadImageToGitHub(t *testing.T) {
	// A valid 1x1 pixel image
	validImg := image.NewRGBA(image.Rect(0, 0, 1, 1))

	// An invalid 0x0 pixel image to intentionally trigger the PNG encoding error
	invalidImg := image.NewRGBA(image.Rect(0, 0, 0, 0))

	cfg := types.UploadStruct{
		Token:     "fake-token",
		Owner:     "testowner",
		Repo:      "testrepo",
		Path:      "stats.png",
		Branch:    "main",
		CommitMsg: "Test Commit",
	}

	tests := []struct {
		name          string
		img           image.Image
		mockRoundTrip func(req *http.Request) (*http.Response, error)
		expectError   bool
		errorContains string
	}{
		{
			name: "PNG Encoding Error",
			img:  invalidImg,
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				return nil, nil // Never reached because encoding fails first
			},
			expectError:   true,
			errorContains: "failed to encode PNG",
		},
		{
			name: "Create Success (File does not exist)",
			img:  validImg,
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					// Simulate file not found (Triggers creation logic)
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Not Found"}`)),
					}, nil
				}
				if req.Method == http.MethodPut {
					// Simulate successful creation
					return &http.Response{
						StatusCode: http.StatusCreated,
						Body:       io.NopCloser(bytes.NewBufferString(`{"content": {"name": "stats.png"}}`)),
					}, nil
				}
				return nil, fmt.Errorf("unexpected request")
			},
			expectError: false,
		},
		{
			name: "Update Success (File already exists)",
			img:  validImg,
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					// Simulate file exists (Provides the SHA needed for an update)
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{"sha": "dummy-sha"}`)),
					}, nil
				}
				if req.Method == http.MethodPut {
					// Simulate successful update
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{"content": {"name": "stats.png"}}`)),
					}, nil
				}
				return nil, fmt.Errorf("unexpected request")
			},
			expectError: false,
		},
		{
			name: "Create Error (API returns 500)",
			img:  validImg,
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Not Found"}`)),
					}, nil
				}
				if req.Method == http.MethodPut {
					// Simulate GitHub API failure during creation
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Internal Server Error"}`)),
					}, nil
				}
				return nil, fmt.Errorf("unexpected request")
			},
			expectError:   true,
			errorContains: "creation failed",
		},
		{
			name: "Update Error (API returns 500)",
			img:  validImg,
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{"sha": "dummy-sha"}`)),
					}, nil
				}
				if req.Method == http.MethodPut {
					// Simulate GitHub API failure during update
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Internal Server Error"}`)),
					}, nil
				}
				return nil, fmt.Errorf("unexpected request")
			},
			expectError:   true,
			errorContains: "update failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup a mock HTTP Client
			mockClient := &http.Client{
				Transport: &mockTransport{roundTripFunc: tt.mockRoundTrip},
			}

			// Inject the mock client into the context.
			// The oauth2 library checks the context for this specific key and uses our mock
			// instead of making a real network request.
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockClient)

			err := UploadImageToGitHub(ctx, tt.img, cfg)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected an error containing %q, but got nil", tt.errorContains)
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain %q, but got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Did not expect an error, but got: %v", err)
				}
			}
		})
	}
}
