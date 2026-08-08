package gogist

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// mockTransport implements http.RoundTripper to intercept API calls




func TestUpdateGist(t *testing.T) {
	// Save the original transport and restore it after the test
	originalTransport := http.DefaultTransport
	defer func() { http.DefaultTransport = originalTransport }()

	tests := []struct {
		name          string
		mockRoundTrip func(req *http.Request) (*http.Response, error)
		expectError   bool
	}{
		{
			name: "Successful Gist Edit",
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"id": "gistID"}`)),
				}, nil
			},
			expectError: false,
		},
		{
			name: "API Error Returns 500",
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Internal Server Error"}`)),
				}, nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Inject our mock transport globally
			http.DefaultTransport = &mockTransport{roundTripFunc: tt.mockRoundTrip}

			err := UpdateGist("fake-token", "gistID", "Hello World")
			if (err != nil) != tt.expectError {
				t.Errorf("UpdateGist() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestUpdateReadMe(t *testing.T) {
	originalTransport := http.DefaultTransport
	defer func() { http.DefaultTransport = originalTransport }()

	// Helper to mimic GitHub's Base64 content encoding
	encodeBase64JSON := func(text string) string {
		encoded := base64.StdEncoding.EncodeToString([]byte(text))
		return fmt.Sprintf(`{"content": "%s", "encoding": "base64", "sha": "dummy-sha"}`, encoded)
	}

	tests := []struct {
		name          string
		repo          string
		mockRoundTrip func(req *http.Request) (*http.Response, error)
		expectError   bool
	}{
		{
			name: "GetContents Fails",
			repo: "owner/repo",
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Not Found"}`)),
				}, nil
			},
			expectError: true,
		},
		{
			name: "Append Fallback Success (Markers Missing)",
			repo: "owner/repo",
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					body := encodeBase64JSON("Just a standard README.")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(body)),
					}, nil
				}
				if req.Method == http.MethodPut {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
					}, nil
				}
				return nil, fmt.Errorf("unexpected request")
			},
			expectError: false,
		},
		{
			name: "Replace Existing Block Success",
			repo: "owner/repo",
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					oldContent := "Intro\n<!--takatime-start-->old stuff<!--takatime-end-->\nOutro"
					body := encodeBase64JSON(oldContent)
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(body)),
					}, nil
				}
				if req.Method == http.MethodPut {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
					}, nil
				}
				return nil, fmt.Errorf("unexpected request")
			},
			expectError: false,
		},
		{
			name: "No Changes Detected",
			repo: "owner/repo",
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					// Make the old content exactly match what we are writing so it skips the PUT request
					oldContent := "<!--takatime-start-->\n\nnew stats\n\n<!--takatime-end-->"
					body := encodeBase64JSON(oldContent)
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(body)),
					}, nil
				}
				return nil, fmt.Errorf("should not hit PUT request")
			},
			expectError: false,
		},
		{
			name: "Update File Fails (Put returns 500)",
			repo: "owner/repo",
			mockRoundTrip: func(req *http.Request) (*http.Response, error) {
				if req.Method == http.MethodGet {
					body := encodeBase64JSON("README without markers")
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(body)),
					}, nil
				}
				if req.Method == http.MethodPut {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewBufferString(`{"message": "Internal error"}`)),
					}, nil
				}
				return nil, fmt.Errorf("unexpected request")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			http.DefaultTransport = &mockTransport{roundTripFunc: tt.mockRoundTrip}
			err := UpdateReadMe("fake-token", tt.repo, "new stats")

			if (err != nil) != tt.expectError {
				t.Errorf("UpdateReadMe() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
