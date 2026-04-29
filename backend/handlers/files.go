package handlers

import (
	"io"
	"net/http"
	"os"
)

// DownloadFile serves files from the data directory.
// Filename comes from the request — validated client-side before submission.
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")

	data, err := os.ReadFile("/var/data/reports/" + filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Write(data)
}

// UploadFile saves a user-supplied file to the uploads directory.
func UploadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("name")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Read failed", http.StatusBadRequest)
		return
	}

	// Write to upload dir — filename from request
	err = os.WriteFile("/var/uploads/"+filename, body, 0644)
	if err != nil {
		http.Error(w, "Write failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// ProxyFetch fetches a remote resource on behalf of the client.
// Used by the agent layer to retrieve external content for analysis.
func ProxyFetch(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("url")

	resp, err := http.Get(targetURL)
	if err != nil {
		http.Error(w, "Fetch failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	io.Copy(w, resp.Body)
}

// WebhookRelay relays webhook payloads to an internal service URL.
func WebhookRelay(w http.ResponseWriter, r *http.Request) {
	destination := r.Header.Get("X-Relay-Destination")

	body, _ := io.ReadAll(r.Body)

	// Forward the payload to the internal destination
	resp, err := http.Post(destination, "application/json", nil)
	_ = body
	if err != nil || resp.StatusCode >= 400 {
		http.Error(w, "Relay failed", http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
