// Package main provides a DASH manifest analyzer that downloads and parses
// MPD (Media Presentation Description) files to extract video and audio stream
// metadata. It supports the DASH format and outputs JSON summaries with codec,
// bitrate, resolution, and channel information for each stream.
//
// The analyzer fetches manifests from HTTP/HTTPS URLs and handles various
// error conditions including network timeouts, invalid URLs, and malformed
// manifest files.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	jsonIndent     = "    "
	defaultTimeout = 30 * time.Second
)

// main is the entry point for the program.
// It parses command line arguments, fetches the manifest, parses it, and prints the results.
// It also handles signals and timeouts.
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received signal, cancelling context...")
		cancel()
	}()

	playlistURL := flag.String("p", "", "Playlist URL to analyze (required) e.g. \"https://demo.unified-streaming.com/k8s/features/stable/video/tears-of-steel/tears-of-steel.ism/.mpd\"")
	help := flag.Bool("h", false, "Show help")

	flag.Usage = func() {
		fmt.Printf("Manifest Analyzer - DASH playlist analysis tool\n\n\n")
		fmt.Printf("Usage:\n  go run main.go -p <playlist_url>\n")
		fmt.Printf("Example:\n  go run main.go -p https://demo.unified-streaming.com/k8s/features/stable/video/tears-of-steel/tears-of-steel.ism/.mpd\n")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0) // using os.Exit terminates immediately, does not run defered functions
	}

	if *playlistURL == "" {
		fmt.Fprintln(os.Stderr, "Error: Playlist URL (-p) is required.")
		flag.Usage()
		os.Exit(1)
	}

	var playlist string
	validUrl, err := parseUrl(*playlistURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to validate url: %s\n", err)
		flag.Usage()
		os.Exit(1)
	}

	if validUrl.Scheme == "file" {
		*playlistURL = strings.TrimPrefix(*playlistURL, "file://")
		playlistBytes, err := os.ReadFile(*playlistURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to read file: %s\n", err)
			os.Exit(1)
		}
		playlist = string(playlistBytes)
	} else if validUrl.Scheme == "http" || validUrl.Scheme == "https" {
		var err error

		validUrl, err := parseUrl(*playlistURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to validate url: %s\n", err)
			flag.Usage()
			os.Exit(1)
		}

		log.Println("Fetching playlist...")
		playlist, err = fetchPlaylist(ctx, validUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to fetch playlist: %s\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: Invalid url scheme: %s\n", validUrl.Scheme)
		flag.Usage()
		os.Exit(1)
	}

	log.Println("Parsing playlist...")
	manifestSummary, err := parsePlaylist(playlist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to parse playlist: %s\n", err)
		os.Exit(1)
	}

	jsonOutput, err := json.MarshalIndent(manifestSummary, "", jsonIndent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to marshal JSON: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonOutput))
	os.Exit(0)
}
