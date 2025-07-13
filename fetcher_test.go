package main

import (
	"context"
	"net/url"
	"testing"
	"time"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid HTTPS URL",
			input:   "https://example.com/manifest.mpd",
			wantErr: false,
		},
		{
			name:    "valid HTTP URL",
			input:   "http://example.com/manifest.mpd",
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			input:   "ftp://example.com/manifest.mpd",
			wantErr: true,
		},
		{
			name:    "empty URL",
			input:   "",
			wantErr: true,
		},
		{
			name:    "malformed URL",
			input:   "not-a-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseUrl(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUrl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFetchPlaylist_InvalidURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invalidURL, _ := url.Parse("https://invalid.com/manifest.mpd")

	_, err := fetchPlaylist(ctx, invalidURL)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}
