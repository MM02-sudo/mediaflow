package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MM02-sudo/mediaflow/shared"
)

func main() {
	// this part tells the server when someone sends a request to /
	http.HandleFunc("/", requestHandler)
	http.HandleFunc("/stream", actualStreamHandler)

	// servers listens for requests on port 8080 it uses default settings(nil) and if it crashes it shows an error and exits with log.Fatal
	fmt.Println("Server starting on port: 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// request handler here
func requestHandler(w http.ResponseWriter, request *http.Request) {
	// reading Json file
	var req shared.Request
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		// if Json file broken it sends an error
		sendError(w, "Invalid JSON")
		return
	}

	// lets check what action the client is requesting
	if req.Action == "list" {
		handleList(w, req.Path)
	} else if req.Action == "stream" {
		handleStream(w, req.Path)
	} else {
		sendError(w, "Unknown action, Try following actions")
	}
}

func handleList(w http.ResponseWriter, path string) {
	// does the path exists
	info, err := os.Stat(path)
	if err != nil {
		sendError(w, "Path not found")
		return
	}

	// lets make sure the path is a directory and not a file
	if !info.IsDir() {
		sendError(w, "Path is not a directory")
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		sendError(w, "Cannot read directory")
		return
	}

	// convering to FileInfo format
	var files []shared.FileInfo
	for _, entry := range entries {
		files = append(files, shared.FileInfo{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		})
	}

	// now send response
	response := shared.Response{
		Success: true,
		Files:   files,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendError(w http.ResponseWriter, message string) {
	response := shared.Response{
		Success: false,
		Error:   message,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleStream(w http.ResponseWriter, path string) {
	// does path exists?
	info, err := os.Stat(path)
	if err != nil {
		sendError(w, "File not found")
		return
	}

	// this time it should be a file not a directory
	if info.IsDir() {
		sendError(w, "Cannot stream directory")
		return
	}

	// Send back the stream URL (don't actually stream yet!)
	streamURL := fmt.Sprintf("http://localhost:8080/stream?path=%s", path)
	response := shared.Response{
		Success:   true,
		StreamURL: streamURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func actualStreamHandler(w http.ResponseWriter, r *http.Request) {
	// Get the file path from URL parameter
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "No path provided", http.StatusBadRequest)
		return
	}

	// Open the video file
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "Cannot open file", http.StatusNotFound)
		return
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		http.Error(w, "Cannot read file info", http.StatusInternalServerError)
		return
	}

	// Set headers for video streaming
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Accept-Ranges", "bytes")

	// Stream file to client
	http.ServeContent(w, r, info.Name(), info.ModTime(), file)
}
