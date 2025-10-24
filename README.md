# Speech to Text Transcription

A beautiful web application for converting speech to text using Deepgram's AI-powered transcription API. Features include speaker diarization, sentiment analysis, topic detection, entity extraction, and more.

## Features

- 🎙️ Audio file upload (MP3, WAV, WebM, etc.)
- 🌐 URL-based transcription
- 👥 Speaker diarization
- 😊 Sentiment analysis
- 🏷️ Entity detection
- 🔖 Topic extraction
- 🎯 Intent detection
- 📝 Smart formatting
- 💬 Paragraph segmentation

## Prerequisites

- Go 1.19 or higher
- Deepgram API key

## Setup

1. Clone the repository:
```bash
git clone https://github.com/MelloB1989/mini-proj-aiml-b
cd mini-proj-aiml-b
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:

Create a `.env` file or export the variable:
```bash
export DEEPGRAM_API_KEY="your_deepgram_api_key_here"
```

To get a Deepgram API key:
- Sign up at https://console.deepgram.com/signup
- Navigate to API Keys section
- Create a new API key

4. Build the application:
```bash
go build -o transcribe main.go
```

## Usage

1. Start the server:
```bash
./transcribe
```

The server will start on `http://localhost:8080`

2. Open your browser and navigate to:
```
http://localhost:8080
```

3. Upload an audio file or paste an audio URL

4. Click "Transcribe" and wait for the results

## API Endpoints

### Upload Audio File
```
POST /api/upload
Content-Type: multipart/form-data

Body: audio file
```

Response:
```json
{
  "success": true,
  "data": {
    "url": "temp/file-path",
    "filename": "unique-filename.mp3"
  }
}
```

### Transcribe Audio
```
POST /api/transcribe
Content-Type: application/json

Body:
{
  "url": "file-path-or-url"
}
```

Response:
```json
{
  "success": true,
  "data": {
    "metadata": {...},
    "results": {
      "channels": [...],
      "summary": {...},
      "topics": {...},
      "intents": {...},
      "sentiments": {...}
    }
  }
}
```

## Project Structure

```
mini-proj/
├── main.go           # Entry point
├── api/
│   └── handler.go    # API handlers
├── frontend/
│   └── index.html    # Web interface
├── temp/             # Temporary uploads (auto-created)
├── go.mod
├── go.sum
└── README.md
```

## Configuration

Default configuration:
- Port: 8080
- Model: nova-3
- Language: English (en)
- Features: All enabled (summarization, topics, intents, entities, sentiment, diarization)

To change settings, modify the options in `api/handler.go`:
```go
options := &interfaces.PreRecordedTranscriptionOptions{
    Model:          "nova-3",
    Language:       "en",
    Summarize:      "v2",
    Topics:         true,
    Intents:        true,
    DetectEntities: true,
    Sentiment:      true,
    SmartFormat:    true,
    Diarize:        true,
    FillerWords:    true,
}
```

## Supported Audio Formats

- MP3
- WAV
- FLAC
- WebM
- OGG
- M4A
- And more...

## Troubleshooting

### "DEEPGRAM_API_KEY not found"
Make sure you've set the environment variable correctly:
```bash
export DEEPGRAM_API_KEY="your_key_here"
```

### Upload fails
- Check file size (max 100MB)
- Ensure the audio format is supported
- Verify temp directory has write permissions

### Transcription fails
- Verify your Deepgram API key is valid
- Check your API usage quota
- Ensure the audio file is accessible
- Check audio quality and format

## Development

To run in development mode:
```bash
go run main.go
```

To modify the port:
Edit `main.go`:
```go
api.StartServer("8080") // Change to your preferred port
```

## License

MIT

## Credits

Powered by [Deepgram](https://deepgram.com/) API
