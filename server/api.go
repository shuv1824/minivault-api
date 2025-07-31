package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var OLLAMA_URL = os.Getenv("OLLAMA_URL")

type APIRequest struct {
	Prompt string `json:"prompt"`
}

type APIResponse struct {
	Response string `json:"response"`
}

type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func Run() {
	http.HandleFunc("/generate", generateHandler)
	log.Println("Go API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	var req APIRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	body, _ := json.Marshal(map[string]any{
		"model":  "llama2",
		"prompt": req.Prompt,
		"stream": true,
	})

	resp, err := http.Post(fmt.Sprintf("http://%s/api/generate", OLLAMA_URL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "Failed to contact Ollama", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the stream line-by-line
	reader := bufio.NewReader(resp.Body)
	var fullResponse strings.Builder

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			http.Error(w, "Failed to read ollama response", http.StatusInternalServerError)
			return
		}

		// Trim possible whitespace
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		var msg OllamaResponse
		if err := json.Unmarshal(line, &msg); err != nil {
			http.Error(w, "Failed to unmarshal ollama response", http.StatusInternalServerError)
			return
		}

		fullResponse.WriteString(msg.Response)

		if msg.Done {
			break
		}
	}

	resBody, _ := json.Marshal(APIResponse{
		Response: fullResponse.String(),
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(resBody)

}
