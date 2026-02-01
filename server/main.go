package main

import (
	"net/http"
	"encoding/json"
	"os"
	"path/filepath"
	"github.com/MM02-sudo/mediaflow/shared"
	"fmt"
	"log"

)


func main(){
	// this part tells the server when someone visits /
	http.HandlerFunc("/", requestHandler)

	// servers listens for requests on port 8080 it uses default settings(nil) and if it crashes it shows an error and exits with log.Fatal
	
	fmt.Println("Server starting on port: 8080...")
	log.Fatal(http.ListenAndServe("8080", nil))

}

//request handler here

func requestHandler(w http.ResponseWriter, request *http.Request)
{

}
