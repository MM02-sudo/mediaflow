package shared

// Request from client to server
type Request struct {
	Action string `json:"action"` // "list" or "stream"
	Path   string `json:"path"`   // file/folder path
}

// Response from server to client
type Response struct {
	Success   bool       `json:"success"`
	Error     string     `json:"error,omitempty"`
	Files     []FileInfo `json:"files,omitempty"`
	StreamURL string     `json:"stream_url,omitempty"`
}

// FileInfo describes a file or folder
type FileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}
