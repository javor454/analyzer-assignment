# Manifest Analyzer - Development Plan

## Overview
Create a command-line tool in Go that analyzes DASH playlists and provides JSON summary output with video and audio stream information.

## Requirements Summary
- Accept `-p` parameter with playlist URI
- Fetch content from URI
- Parse DASH manifest
- Output JSON with video/audio stream summaries
- Video: codec, resolution, bitrate
- Audio: codec, number of channels, language

## Phase 1: Project Setup & Architecture Design

### 1.1 Simplified Project Structure ✅ **COMPLETED**
```
manifest-analyzer-assignment/
├── main.go              # Single file containing all components
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── development.md
```

**Implementation Status:**
- ✅ All components will be implemented in main.go
- ✅ Single file approach as required by README
- ✅ Ready for component implementation within main.go

### 1.2 Technology Stack Decisions ✅ **COMPLETED**
- **DASH Support**: Use `github.com/zencoder/go-dash` (already in go.mod)
- **HTTP Client**: Standard `net/http` with timeout handling
- **CLI Framework**: `flag` package (standard library) ✅ **IMPLEMENTED**
- **JSON Output**: Standard `encoding/json`
- **Error Handling**: Custom error types with context

**Implementation Status:**
- ✅ CLI Framework implemented using `flag` package
- ✅ DASH library already in go.mod
- ✅ Ready for HTTP client implementation
- ✅ Ready for JSON output implementation

## Phase 2: Core Components Development

### 2.1 CLI Interface ✅ **COMPLETED**
- Parse `-p` flag for playlist URI
- Validate URI format
- Handle help/usage display
- Set up proper error handling and exit codes

**Implementation Status:**
- ✅ Flag parsing with `-p` and `-h` flags
- ✅ Custom usage function with clear examples
- ✅ URL validation using `url.ParseRequestURI()`
- ✅ Proper error handling with `os.Exit()` and clear error messages
- ✅ Good separation of concerns (validation vs usage display)

### 2.2 Content Fetcher (in `main.go`)
- HTTP client with timeout (30s default)
- User-Agent header to avoid blocking
- Handle redirects
- Validate content type (text/plain, application/xml, etc.)
- Error handling for network issues

#### Content Fetcher Plan Details

**Core Responsibilities:**

1. **HTTP Client Setup**
   - Create configured HTTP client with appropriate timeouts
   - Set reasonable timeout values (30 seconds default)
   - Configure user-agent header to avoid being blocked by servers
   - Handle redirects automatically (up to a reasonable limit)

2. **URL Validation & Processing**
   - Validate that the URL is properly formatted
   - Ensure it's an HTTP/HTTPS URL
   - Handle relative URLs if needed
   - Sanitize the URL to prevent injection attacks

3. **Content Fetching**
   - Make HTTP GET request to the playlist URL
   - Handle different HTTP status codes appropriately:
     - 200: Success
     - 301/302: Redirects (handled automatically)
     - 404: Not found
     - 403: Forbidden
     - 500+: Server errors
   - Validate content type (expecting XML for DASH manifests)
   - Handle large files efficiently (streaming vs buffering)

4. **Error Handling**
   - Network timeouts
   - DNS resolution failures
   - SSL/TLS certificate issues
   - Malformed URLs
   - Server errors
   - Empty responses
   - Invalid content types

5. **Response Processing**
   - Read the response body
   - Handle different content encodings (gzip, deflate)
   - Validate that we received actual content (not empty)
   - Check content length limits (prevent memory issues)
   - Return the raw manifest content as string

**Key Design Decisions:**

- **Timeout Strategy**: Connection timeout (10s), Read timeout (30s), Total timeout (60s max)
- **User-Agent Header**: Set proper user-agent to avoid being blocked
- **Content Validation**: Check Content-Type header, validate response size (max 10MB)
- **Error Types**: Create custom error types for different failure scenarios

**Function Signature Plan:**
```go
// Main function to fetch manifest content
func fetchManifest(url string) (string, error)

// Helper function to validate URL
func validateURL(url string) error

// Helper function to create HTTP client
func createHTTPClient() *http.Client

// Helper function to validate response
func validateResponse(resp *http.Response) error
```

### 2.3 Parser Components (in `main.go`)
- **DASH Parser**: Parse MPD (Media Presentation Description) XML
- Extract adaptation sets and representations
- Parse codec strings and media info

### 2.4 Analyzer Engine (in `main.go`)
- Process parsed manifest data
- Extract video/audio stream information
- Calculate bitrates and resolutions
- Handle codec parsing and validation

### 2.5 Data Models (in `main.go`)
```go
type ManifestSummary struct {
    Videos []VideoStream `json:"videos"`
    Audios []AudioStream `json:"audios"`
}

type VideoStream struct {
    Codec      string `json:"codec"`
    Bitrate    string `json:"bitrate"`
    Resolution string `json:"resolution"`
}

type AudioStream struct {
    Codec    string `json:"codec"`
    Bitrate  string `json:"bitrate"`
    Channels int    `json:"channels"`
    Language string `json:"language"`
}
```

## Phase 3: Implementation Strategy

### 3.1 Development Order
1. **CLI Framework** - Basic flag parsing and help
2. **HTTP Fetcher** - Test with simple URLs first
3. **DASH Parser** - Start with basic XML parsing
4. **Data Models** - Define structures
5. **Analyzer Logic** - Extract stream information
6. **JSON Output** - Format results
7. **Error Handling** - Comprehensive error management
8. **Testing** - Unit tests for each component

### 3.2 Key Implementation Challenges
- **Codec Parsing**: Handle various codec formats (avc1.42C00D, mp4a.40.2)
- **Bitrate Calculation**: Extract from bandwidth attributes
- **Resolution Parsing**: Handle different resolution formats
- **Language Detection**: Extract from adaptation set attributes
- **Error Recovery**: Handle malformed manifests gracefully

## Phase 4: Testing & Validation

### 4.1 Test Strategy
- **Unit Tests**: Each component in isolation
- **Integration Tests**: End-to-end with sample manifests
- **Error Tests**: Invalid URLs, malformed manifests
- **Performance Tests**: Large manifest files

### 4.2 Sample Test Cases
- Valid DASH manifest (provided example)
- Invalid URI handling
- Network timeout scenarios
- Malformed XML handling
- Empty manifest handling

## Phase 5: Enhancement & Polish

### 5.1 Optional Features
- **HLS Support**: Add M3U8 parsing capability
- **Verbose Mode**: Detailed parsing information
- **Output Format**: Support for different output formats
- **Caching**: Cache fetched manifests
- **Progress Indicators**: For large manifests

### 5.2 Code Quality
- **Documentation**: Comprehensive godoc comments
- **Logging**: Structured logging for debugging
- **Configuration**: Environment-based settings
- **Performance**: Optimize for large manifests

## Phase 6: Deployment & Distribution

### 6.1 Build Process
- Cross-platform builds (Linux, macOS, Windows)
- Docker containerization
- CI/CD pipeline setup

### 6.2 Distribution
- Binary releases
- Docker images
- Installation scripts

## Risk Assessment & Mitigation

### High Risk Areas:
1. **Manifest Format Variations**: Different DASH implementations
2. **Network Reliability**: Timeout and retry strategies
3. **Codec Complexity**: Various codec string formats
4. **Performance**: Large manifest files

### Mitigation Strategies:
1. **Robust Parsing**: Handle edge cases gracefully
2. **Comprehensive Testing**: Multiple manifest formats
3. **Error Recovery**: Graceful degradation
4. **Performance Monitoring**: Profile and optimize

## Success Criteria
- ✅ Parses provided example manifest correctly
- ✅ Generates exact JSON output as specified
- ✅ Handles network errors gracefully
- ✅ Provides clear error messages
- ✅ Fast execution (< 5 seconds for typical manifests)
- ✅ Memory efficient (< 50MB for large manifests)

## Example Output Format
```json
{
    "audios": [
        {
            "codec": "mp4a.40.2",
            "bitrate": "64008",
            "channels": 2,
            "language": "en"
        },
        {
            "codec": "mp4a.40.2",
            "bitrate": "128008",
            "channels": 2,
            "language": "en"
        }
    ],
    "videos": [
        {
            "codec": "avc1.42C00D",
            "bitrate": "401000",
            "resolution": "224x100"
        },
        {
            "codec": "avc1.42C016",
            "bitrate": "751000",
            "resolution": "448x200"
        },
        {
            "codec": "avc1.4D401F",
            "bitrate": "1001000",
            "resolution": "784x350"
        },
        {
            "codec": "avc1.640028",
            "bitrate": "1501000",
            "resolution": "1680x750"
        },
        {
            "codec": "avc1.640028",
            "bitrate": "2200000",
            "resolution": "1680x750"
        }
    ]
}
```

## Command Usage
```bash
go run main.go -p https://demo.unified-streaming.com/k8s/features/stable/video/tears-of-steel/tears-of-steel.ism/.mpd
``` 