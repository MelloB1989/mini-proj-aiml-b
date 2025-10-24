package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/MelloB1989/karma/config"
	"github.com/MelloB1989/karma/utils"
	listen "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
)

type TranscribeResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func EnableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func UploadMedia(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		json.NewEncoder(w).Encode(TranscribeResponse{
			Success: false,
			Error:   "Failed to parse form",
		})
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		json.NewEncoder(w).Encode(TranscribeResponse{
			Success: false,
			Error:   "Failed to get file",
		})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", utils.GenerateID(6), ext)
	tempDir := "./temp"
	os.MkdirAll(tempDir, os.ModePerm)

	filepath := filepath.Join(tempDir, filename)
	dst, err := os.Create(filepath)
	if err != nil {
		json.NewEncoder(w).Encode(TranscribeResponse{
			Success: false,
			Error:   "Failed to create file",
		})
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		json.NewEncoder(w).Encode(TranscribeResponse{
			Success: false,
			Error:   "Failed to save file",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TranscribeResponse{
		Success: true,
		Data: map[string]string{
			"url":      filepath,
			"filename": filename,
		},
	})
}

func Transcribe(w http.ResponseWriter, r *http.Request) {
	EnableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		json.NewEncoder(w).Encode(TranscribeResponse{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	client.InitWithDefault()
	ctx := context.Background()

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

	c := client.NewREST(config.GetEnvRaw("DEEPGRAM_API_KEY"), &interfaces.ClientOptions{})
	dg := listen.New(c)

	var res any
	if _, err := os.Stat(req.URL); err == nil {
		file, err := os.Open(req.URL)
		if err != nil {
			json.NewEncoder(w).Encode(TranscribeResponse{
				Success: false,
				Error:   "Failed to open file",
			})
			return
		}
		defer file.Close()
		defer os.Remove(req.URL)

		res, err = dg.FromStream(ctx, file, options)
		if err != nil {
			json.NewEncoder(w).Encode(TranscribeResponse{
				Success: false,
				Error:   fmt.Sprintf("Transcription failed: %v", err),
			})
			return
		}
	} else {
		res, err = dg.FromURL(ctx, req.URL, options)
		if err != nil {
			json.NewEncoder(w).Encode(TranscribeResponse{
				Success: false,
				Error:   fmt.Sprintf("Transcription failed: %v", err),
			})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TranscribeResponse{
		Success: true,
		Data:    res,
	})
}

func StartServer(port string) {
	http.HandleFunc("/api/upload", UploadMedia)
	http.HandleFunc("/api/transcribe", Transcribe)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
		os.Exit(1)
	}
}
