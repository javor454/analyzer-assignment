package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zencoder/go-dash/mpd"
)

const (
	audioContentType = "audio"
	videoContentType = "video"
	audioMimeType    = "audio/mp4"
	videoMimeType    = "video/mp4"
	unknown          = "unknown"
)

type (
	// ManifestSummary is a summary of the manifest, which describes audio and video stream metadata
	ManifestSummary struct {
		Audios []AudioStream `json:"audios,omitempty"`
		Videos []VideoStream `json:"videos,omitempty"`
	}
	VideoStream struct {
		Codec      string `json:"codec"`
		Bitrate    string `json:"bitrate"`
		Resolution string `json:"resolution"`
	}
	AudioStream struct {
		Codec    string `json:"codec"`
		Bitrate  string `json:"bitrate"`
		Channels string `json:"channels"` // hex for dolby, int for others
		Language string `json:"language"`
	}
)

// parsePlaylist parses the given manifest string and returns a ManifestSummary object
func parsePlaylist(manifestBytes string) (ManifestSummary, error) {
	var summary ManifestSummary
	mpd, err := mpd.ReadFromString(manifestBytes) // all lib functions save whole file to memory
	if err != nil {
		return summary, fmt.Errorf("failed to read DASH manifest: %w", err)
	}

	for _, period := range mpd.Periods {
		for i, adaptationSet := range period.AdaptationSets {
			if adaptationSet.ContentType == nil || adaptationSet.MimeType == nil {
				fmt.Printf("AdaptationSet %d: ContentType or MimeType is nil, skipping...\n", i)
				continue
			}

			if *adaptationSet.ContentType == audioContentType && *adaptationSet.MimeType == audioMimeType {
				for _, representation := range adaptationSet.Representations {
					codecs := extractCodec(adaptationSet.Codecs)
					if codecs == "unknown" {
						codecs = extractCodec(representation.Codecs)
					}

					channels := extractChannelDescriptors(adaptationSet.AudioChannelConfiguration)
					if channels == "unknown" {
						channels = extractChannelConfiguration(representation.AudioChannelConfiguration)
					}

					audioStream := AudioStream{
						Codec:    codecs,
						Bitrate:  extractBitrate(representation.Bandwidth),
						Channels: channels,
						Language: extractLanguage(adaptationSet.Lang),
					}
					summary.Audios = append(summary.Audios, audioStream)
				}
			}
			if *adaptationSet.ContentType == videoContentType && *adaptationSet.MimeType == videoMimeType {
				for _, representation := range adaptationSet.Representations {
					codecs := extractCodec(adaptationSet.Codecs)
					if codecs == "unknown" {
						codecs = extractCodec(representation.Codecs)
					}

					videoStream := VideoStream{
						Codec:      codecs,
						Bitrate:    extractBitrate(representation.Bandwidth),
						Resolution: extractResolution(representation.Width, representation.Height),
					}
					summary.Videos = append(summary.Videos, videoStream)
				}
			}
		}
	}

	return summary, nil
}

func extractCodec(codecs *string) string {
	if codecs == nil || *codecs == "" {
		return unknown
	}

	return *codecs
}

func extractResolution(width, height *int64) string {
	if width == nil || height == nil || *width == 0 || *height == 0 {
		return unknown
	}

	return fmt.Sprintf("%dx%d", *width, *height)
}

func extractBitrate(bandwidth *int64) string {
	if bandwidth == nil || *bandwidth == 0 {
		return unknown
	}

	return strconv.FormatInt(*bandwidth, 10)
}

func extractChannelDescriptors(descriptors []mpd.DescriptorType) string {
	channels := make([]string, 0, len(descriptors))
	for _, descriptor := range descriptors {
		if descriptor.Value == nil || *descriptor.Value == "" {
			continue
		}

		channels = append(channels, *descriptor.Value)
	}

	if len(channels) == 0 {
		return unknown
	}

	return strings.Join(channels, ",")
}

func extractChannelConfiguration(channelConfiguration *mpd.AudioChannelConfiguration) string {
	if channelConfiguration == nil || channelConfiguration.Value == nil || *channelConfiguration.Value == "" {
		return unknown
	}

	return *channelConfiguration.Value
}

func extractLanguage(lang *string) string {
	if lang == nil || *lang == "" {
		return unknown
	}

	return *lang
}
