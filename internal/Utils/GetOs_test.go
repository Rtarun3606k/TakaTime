package utils

import (
	"runtime"
	"testing"
)

func TestGetOS(t *testing.T) {
	got := GetOS()
	want := runtime.GOOS

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestGetOS_NotEmpty(t *testing.T) {
	got := GetOS()

	if got == "" {
		t.Fatal("expected non-empty OS name")
	}
}
