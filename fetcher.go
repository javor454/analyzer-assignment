package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	mpegDashMPDMimeType  = "application/dash+xml"
	appXmlMimeType       = "application/xml"
	textXmlMimeType      = "text/xml"
	maxPlaylistSizeBytes = 10 * 1024 * 1024 // 10MB
	userAgent            = "ManifestAnalyzer/1.0 (GO HTTP Client)"
)

// fetchPlaylist fetches the playlist from the given url and returns the playlist as a string
// it supports mpeg dash, application/xml, and text/xml mime types
func fetchPlaylist(ctx context.Context, url *url.URL) (string, error) {
	client := &http.Client{Timeout: defaultTimeout}

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent) // avoid being blocked
	req.Header.Set("Accept", fmt.Sprintf("%s, %s, %s", mpegDashMPDMimeType, appXmlMimeType, textXmlMimeType))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got non-200 status code: %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != mpegDashMPDMimeType && resp.Header.Get("Content-Type") != appXmlMimeType && resp.Header.Get("Content-Type") != textXmlMimeType {
		return "", fmt.Errorf("got unexpected content type: %s", resp.Header.Get("Content-Type"))
	}

	if resp.ContentLength == 0 {
		return "", fmt.Errorf("got empty response body")
	}

	if resp.ContentLength > maxPlaylistSizeBytes {
		return "", fmt.Errorf("playlist is too large: %d bytes", resp.ContentLength)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read all bytes: %w", err)
	}

	return string(bytes), nil
}

// parseUrl parses the given url and returns a valid url.URL object
func parseUrl(rawUrl string) (*url.URL, error) {
	validUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	if validUrl.Scheme != "http" && validUrl.Scheme != "https" && validUrl.Scheme != "file" {
		return nil, fmt.Errorf("invalid url scheme: %s", validUrl.Scheme)
	}

	return validUrl, nil
}
