package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"bufio"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Payload struct {
	Model   string    `json:"model"`
	Messages []Message `json:"messages"`
}

func main() {
	client := &http.Client{}
	scanner := bufio.NewScanner(os.Stdin)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable not set")
		return
	}

	for {
		fmt.Print("You: ")
		scanner.Scan()
		userMessage := scanner.Text()

		if userMessage == "exit" {
			break
		}

		payload := Payload{
			Model: "gpt-3.5-turbo",
			Messages: []Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: userMessage},
			},
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Error marshalling payload:", err)
			return
		}

		req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer " + apiKey)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}

		if choices, exists := response["choices"].([]interface{}); exists {
			if len(choices) > 0 {
				choice := choices[0].(map[string]interface{})
				if message, exists := choice["message"].(map[string]interface{}); exists {
					if content, exists := message["content"].(string); exists {
						fmt.Println("ChatGPT:", content)
					} else {
						fmt.Println("Content field is missing or not a string")
					}
				} else {
					fmt.Println("Message field is missing or not a map")
				}
			} else {
				fmt.Println("No response choices received")
			}
		} else {
			fmt.Println("Unexpected response format")
		}
	}
}
