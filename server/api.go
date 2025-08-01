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
	"runtime"
	"strings"
	"time"
)

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

type Status struct {
	Uptime       string  `json:"uptime"`
	MemoryUsedMB float64 `json:"memory_used_mb"`
	NumGoroutine int     `json:"num_goroutine"`
	NumCPU       int     `json:"num_cpu"`
}

func Run(startTime time.Time) {
	http.HandleFunc("/generate", generateHandler)        // POST method
	http.HandleFunc("/status", statusHandler(startTime)) // GET method

	log.Println("Go API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05") // standard Go timestamp format

	var req APIRequest

	ollama_url := os.Getenv("OLLAMA_URL")
	model := os.Getenv("MODEL")
	if model == "" {
		model = "llama2"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	body, _ := json.Marshal(map[string]any{
		"model":  model,
		"prompt": req.Prompt,
		"stream": true,
	})

	resp, err := http.Post(fmt.Sprintf("%s/api/generate", ollama_url), "application/json", bytes.NewBuffer(body))
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

	logToFile(timestamp, req.Prompt, fullResponse.String())

	w.Header().Set("Content-Type", "application/json")
	w.Write(resBody)

}

func statusHandler(startTime time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		status := Status{
			Uptime:       time.Since(startTime).String(),
			MemoryUsedMB: float64(memStats.Alloc) / 1024 / 1024,
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

func logToFile(timestamp, prompt, response string) {
	f, err := os.OpenFile("logs/chat.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	logEntry := fmt.Sprintf("[%s]\nPrompt: %s\nResponse: %s\n\n", timestamp, prompt, response)

	f.WriteString(logEntry)
}
