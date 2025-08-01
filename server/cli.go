package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func RunCLI() {
	ollama_url := os.Getenv("OLLAMA_URL")
	model := os.Getenv("MODEL")
	if model == "" {
		model = "llama2"
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Running CLI against model: %s (type 'exit' to quit)\n\n", model)

	for {
		fmt.Print("Prompt: ")
		prompt, _ := reader.ReadString('\n')
		prompt = strings.TrimSpace(prompt)

		if prompt == "" {
			continue
		}
		if prompt == "exit" || prompt == "quit" {
			fmt.Println("Exiting...")
			break
		}

		body, _ := json.Marshal(map[string]any{
			"model":  model,
			"prompt": prompt,
			"stream": true,
		})

		resp, err := http.Post(fmt.Sprintf("%s/api/generate", ollama_url), "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("Error connecting to Ollama:", err)
			continue
		}
		defer resp.Body.Close()

		fmt.Print("=> ")

		// Ollama streams newline-delimited JSON lines (NDJSON)
		decoder := json.NewDecoder(resp.Body)
		for decoder.More() {
			var chunk map[string]any
			if err := decoder.Decode(&chunk); err != nil {
				break
			}
			if token, ok := chunk["response"].(string); ok {
				fmt.Print(token)
			}
		}
		fmt.Println()
	}
}
