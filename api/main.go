package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var UTCBuildTime string

func pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func statusHandler(w http.ResponseWriter, _ *http.Request) {
	bytes, _ := json.Marshal(struct {
		Version      string `json:"version"`
		Status       string `json:"status"`
		UTCBuildTime string `json:"utc_build_time"`
		BucketName   string `json:"bucket_name"`
	}{
		Version:      "0.0.6",
		Status:       "ok",
		UTCBuildTime: UTCBuildTime,
		BucketName:   os.Getenv("DOCUMENTS_BUCKET_NAME"),
	})

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/status", statusHandler)
	fmt.Println("Server listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
