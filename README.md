# Simple Manifest Analyzer – Homework - variant A.1
## How to run
- `make build` - build binary
- `make run` - build binary + run with parameter from assignment
- `make test` - run tests
- `make test-coverage` - run tests with coverage + open coverage explorer

## Assignment 
- Create a command-line tool in Go which analyze playlists (choose one of DASH or HLS format to support) and provides summary to stdout
- Accept parameter `–p` with playlist URI
- Fetch content of playlist from URI
- Parse the playlist
- Create JSON output with summary per each respective stream:
    - Video: codec, resolution, bitrate
    - Audio: codec, number of channels, language
- Feel free to fill the gaps in this task with your creativity
- Send us source code (*.go files, go.mod and go.sum) of completed solution in ZIP
file
- Command example:
`go run main.go -p https://demo.unified-streaming.com/k8s/features/stable/video/tears-of-steel/tears-of-steel.ism/.mpd`
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
Recommended libraries
- https://github.com/etherlabsio/go-m3u8
- https://github.com/zencoder/go-dash/

#mpd file structure
```
MPD
└── Period
    ├── AdaptationSet (audio)
    │   ├── AudioChannelConfiguration
    │   ├── Role
    │   ├── SegmentTemplate
    │   │   └── SegmentTimeline
    │   └── Representation (multiple)
    └── AdaptationSet (video)
        ├── Role
        ├── SegmentTemplate
        │   └── SegmentTimeline
        └── Representation (multiple)
```

## TODOs
- dockerize
- tests
    - more test cases
    - integration tests