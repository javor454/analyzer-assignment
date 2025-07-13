package main

import (
	"testing"
)

func TestExtractCodec(t *testing.T) {
	tests := []struct {
		name     string
		codecs   *string
		expected string
	}{
		{
			name:     "valid codec",
			codecs:   stringPtr("avc1.42C00D"),
			expected: "avc1.42C00D",
		},
		{
			name:     "nil codec",
			codecs:   nil,
			expected: "unknown",
		},
		{
			name:     "empty codec",
			codecs:   stringPtr(""),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractCodec(tt.codecs)
			if result != tt.expected {
				t.Errorf("extractCodec() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractResolution(t *testing.T) {
	tests := []struct {
		name     string
		width    *int64
		height   *int64
		expected string
	}{
		{
			name:     "valid resolution",
			width:    int64Ptr(1920),
			height:   int64Ptr(1080),
			expected: "1920x1080",
		},
		{
			name:     "nil dimensions",
			width:    nil,
			height:   nil,
			expected: "unknown",
		},
		{
			name:     "zero dimensions",
			width:    int64Ptr(0),
			height:   int64Ptr(0),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractResolution(tt.width, tt.height)
			if result != tt.expected {
				t.Errorf("extractResolution() = %v, want %v", result, tt.expected)
			}
		})
	}
}
