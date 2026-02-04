package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/MM02-sudo/mediaflow/shared"
)

func main() {
	server := flag.String("server", "", "Server address(e.g., 192.168.1.100:8080)")
	action := flag.String("action", "list", "Action: list or stream")
	path := flag.String("path", "/", "File or Folder path")

	flag.Parse()

	// lets check if server was provided
	if *server == "" {
		fmt.Println("Error: -server is required")
		flag.Usage()
		return
	}

	// creating a request
	req := shared.Request{
		Action: *action,
		Path:   *path,
	}

	// converting request into json
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error while creating request:", err)
		return
	}

	// sending an http POST to server
	url := fmt.Sprintf("http://%s/", *server)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	var response shared.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	displayResponse(response)
}

func displayResponse(resp shared.Response) {
	if !resp.Success {
		fmt.Println("Error:", resp.Error)
		return
	}

	fmt.Println("Success!")

	if len(resp.Files) > 0 {
		fmt.Println("\nFiles and Folders:")
		for _, file := range resp.Files {
			if file.IsDir {
				fmt.Printf("  [DIR]  %s\n", file.Name)
			} else {
				fmt.Printf("  [FILE] %s\n", file.Name)
			}
		}
	}

	if resp.StreamURL != "" {
		fmt.Println("\nStream URL:", resp.StreamURL)
		fmt.Println("Open this URL in VLC or your browser to play the video!")
	}
}
